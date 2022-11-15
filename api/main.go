package api

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

//go:embed data.json
var jsonContent []byte

type API struct {
	counters counters
	data     []User
}

type counters struct {
	mu       sync.Mutex
	errors   int
	requests int
}

func (a *API) GetUsers() ([]User, error) {
	return a.data, nil
}

func (a *API) GetUser(name string) (User, error) {
	for _, user := range a.data {
		if user.Name == name {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("user not found")
}
func (a *API) GetCounters() counters {
	return a.counters
}
func (c *counters) String() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return fmt.Sprintf("Requests: %d, Errors: %d", c.requests, c.errors)
}

func New() (*API, error) {
	content, err := getDatabase()
	if err != nil {
		return nil, fmt.Errorf("db init: %w", err)
	}
	a := &API{
		data:     content,
		counters: counters{},
	}
	return a, nil
}

func (a *API) Run(ctx context.Context, addr string) error {
	router := gin.Default()
	router.GET("/user/:name", a.getUser)
	router.POST("/user/:name/addpet", a.addPet)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
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
