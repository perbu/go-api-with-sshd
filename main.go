package main

import (
	"context"
	"fmt"
	"github.com/perbu/go-api-with-sshd/api"
	"github.com/perbu/go-api-with-sshd/backdoor"
	"log"
	"os"
	"os/signal"
	"sync"
)

func realMain() error {
	a, err := api.New()
	if err != nil {
		return fmt.Errorf("api init: %w", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.Run(ctx, ":8080")
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := backdoor.Run(ctx, ":2222", a)
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	err := realMain()
	if err != nil {
		log.Panicln(err)
	}
}
