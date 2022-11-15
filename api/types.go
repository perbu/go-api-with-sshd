package api

import (
	"fmt"
	"sync"
)

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

func (u *User) AddLog(log string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.logs == nil {
		u.logs = make([]string, 0)
	}
	u.logs = append(u.logs, log)
}
func (u *User) GetLogs() []string {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.logs
}

type Pet struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
