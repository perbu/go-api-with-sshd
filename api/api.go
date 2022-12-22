package api

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
)

//go:embed data.json
var jsonContent []byte

type API struct {
	counters *Counters
	data     []*User
}

type Counters struct {
	mu       sync.Mutex
	errors   int
	requests int
}

func (a *API) GetUsers() ([]UserDTO, error) {
	resp := make([]UserDTO, 0, len(a.data))
	for _, user := range a.data {
		resp = append(resp, user.DTO())
	}
	return resp, nil
}

func (a *API) GetUser(name string) (UserDTO, error) {
	for _, user := range a.data {
		if user.Name == name {
			return user.DTO(), nil
		}
	}
	return UserDTO{}, fmt.Errorf("user not found")
}
func (a *API) GetCounters() *Counters {
	return a.counters
}
func (c *Counters) String() string {
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
		counters: &Counters{},
	}
	return a, nil
}

func getDatabase() ([]*User, error) {
	var db []*User
	err := json.Unmarshal(jsonContent, &db)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal database: %w", err)
	}
	return db, nil
}
