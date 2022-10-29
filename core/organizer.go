package core

import "fmt"

type User struct {
	Name string
	ID   string
}

// Mention generated a Discord mention string for the user
func (user *User) Mention() string {
	return fmt.Sprintf("<@%s>", user.ID)
}
