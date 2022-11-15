// Package backdoor provides a simple embedded ssh server that can be used to
// access the application through a cli
package backdoor

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/perbu/go-api-with-sshd/api"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"strings"
)

//go:embed authorized_keys
var embeddedAuthorizedKeys []byte

type sshApp struct {
	server *ssh.Server
	pubKey ssh.PublicKey
	api    *api.API
}

func getAuthorizedKeys() ssh.PublicKey {
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(embeddedAuthorizedKeys)
	if err != nil {
		log.Fatalf("Couldn't make sense of authorized key: %s", err)
	}
	return pubKey
}

func (a sshApp) sshHandler(s ssh.Session) {
	defer s.Close()
	if s.RawCommand() != "" {
		_, _ = io.WriteString(s, "raw commands are not supported")
		return
	}
	term := terminal.NewTerminal(s, fmt.Sprintf("%s> ", s.User()))
	pty, winCh, isPty := s.Pty()
	if isPty {
		fmt.Println("PTY term", pty.Term)
		go func() {
			for chInfo := range winCh {
				err := term.SetSize(chInfo.Width, chInfo.Height)
				if err != nil {
					fmt.Println("winch error:", err)
				}
			}
		}()
	}
	_, err := io.WriteString(s, fmt.Sprintf("Welcome, %s\n", s.User()))
	if err != nil {
		log.Println(err)
		return
	}

	for {
		line, err := term.ReadLine()
		if err == io.EOF {
			// Ignore errors here:
			_, _ = io.WriteString(s, "EOF.\n")
			break
		}
		if err != nil {
			// Ignore errors here:
			_, _ = io.WriteString(s, "Error while reading: "+err.Error())
			break
		}
		if line == "quit" {
			break
		}
		if line == "" {
			continue
		}
		output, err := a.handleTerminalInput(line)
		if err != nil {
			log.Printf("Error handling terminal input: %s", err)
			return
		}
		_, err = io.WriteString(s, output)
		if err != nil {
			log.Printf("Error writing to session: %s", err)
			return
		}
	}
}

func (a sshApp) handleTerminalInput(line string) (string, error) {
	ss := strings.SplitN(line, " ", 2)
	switch ss[0] {
	case "ls":
		return a.handleLs()
	case "logs":
		if len(ss) < 2 {
			return "logs command requires a user name\n", nil
		}
		return a.handleLogs(ss[1])
	case "echo":
		return fmt.Sprintf("echo: %s\n", line), nil
	default:
		return fmt.Sprintf("command not recognized\n"), nil
	}
}

func (a sshApp) handleLs() (string, error) {
	users, err := a.api.GetUsers()
	if err != nil {
		return "", err
	}
	var output string
	for _, user := range users {
		output += user.String() + "\n"
	}
	return output, nil
}

func (a sshApp) handleLogs(userName string) (string, error) {
	user, err := a.api.GetUser(userName)
	if err != nil {
		return "", err
	}
	logs := user.GetLogs()
	output := fmt.Sprintf("Logs for user %s (%d lines):\n", user.Name, len(logs))
	for _, logline := range logs {
		output += logline + "\n"
	}
	return output, nil
}

func (a sshApp) myPubKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	if ssh.KeysEqual(key, a.pubKey) {
		return true
	} else {
		return false
	}
}

func Run(ctx context.Context, addr string, api *api.API) error {
	a := sshApp{
		pubKey: getAuthorizedKeys(),
		api:    api,
	}
	a.server = &ssh.Server{
		Addr:             addr,
		Handler:          a.sshHandler,
		PublicKeyHandler: a.myPubKeyHandler,
	}
	go func() {
		<-ctx.Done()
		_ = a.server.Shutdown(context.TODO())
	}()
	log.Println("Starting ssh server on", addr)
	return a.server.ListenAndServe()
}
