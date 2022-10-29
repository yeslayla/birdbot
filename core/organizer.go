package core

import "fmt"

type User struct {
	Name string
	ID   string
}

func (user *User) Mention() string {
	return fmt.Sprintf("<@%s>", user.ID)
}
