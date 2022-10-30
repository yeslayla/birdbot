package core

import "fmt"

type User struct {
	ID string
}

// Mention generated a Discord mention string for the user
func (user *User) Mention() string {
	if user == nil {
		return "<NULL>"
	}

	return fmt.Sprintf("<@%s>", user.ID)
}
