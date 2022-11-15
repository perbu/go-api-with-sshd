package api

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//go:embed data.json
var jsonContent []byte

type API struct {
	counter counters
	data    []User
}

type counters struct {
	errors   int
	requests int
}

func New() (*API, error) {
	content, err := getDatabase()
	if err != nil {
		return nil, fmt.Errorf("db init: %w", err)
	}
	a := &API{
		data:    content,
		counter: counters{},
	}
	return a, nil
}

func (a *API) Run(ctx context.Context, addr string) error {
	srv := &http.Server{
		Addr: addr,
	}
	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5)
		defer cancel()
		err := srv.Shutdown(ctxShutdown)
		if err != nil {
			log.Fatalln("Unexpected http shutdown error:", err)
		}
	}()
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server: %w", err)
	}
	return nil
}

func getDatabase() ([]User, error) {
	var db []User
	err := json.Unmarshal(jsonContent, &db)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal database: %w", err)
	}
	return db, nil
}
