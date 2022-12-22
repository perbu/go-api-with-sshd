package api

import (
	"fmt"
	"sync"
	"time"
)

// User is a struct that represents a user in the database.
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Pets []Pet  `json:"pets"`
	logs []string
	mu   sync.Mutex
}

func (u *User) String() string {
	u.mu.Lock()
	defer u.mu.Unlock()
	return fmt.Sprintf("User: %s, %d, %v", u.Name, u.Age, u.Pets)
}

// UserDTO is the DTO for the User struct.
type UserDTO struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Pets []Pet  `json:"pets"`
	logs []string
}

func (u UserDTO) String() string {
	return fmt.Sprintf("User: %s, %d, %v", u.Name, u.Age, u.Pets)
}

func (u *User) DTO() UserDTO {
	u.mu.Lock()
	defer u.mu.Unlock()
	return UserDTO{
		Name: u.Name,
		Age:  u.Age,
		Pets: u.Pets,
		logs: u.logs,
	}
}

func (u *User) AddLog(log string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.logs == nil {
		u.logs = make([]string, 0)
	}
	// add a timestamp to the log, so it looks nice and proper:
	log = fmt.Sprintf("%v: %s", time.Now().Format(time.RFC3339), log)
	u.logs = append(u.logs, log)
}

func (u *User) GetLogs() []string {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.logs
}

func (u UserDTO) GetLogs() []string {
	return u.logs
}

type Pet struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
