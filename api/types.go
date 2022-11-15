package api

import "fmt"

type User struct {
	Name string   `json:"name"`
	Age  int      `json:"age"`
	Pets []string `json:"pets"`
}

func (u User) String() string {
	return fmt.Sprintf("User: %s, %d, %v", u.Name, u.Age, u.Pets)
}
