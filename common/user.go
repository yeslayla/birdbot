package common

import "fmt"

type User struct {
	ID          string
	AvatarURL   string
	DisplayName string
}

// DiscordMention generated a Discord mention string for the user
func (user *User) DiscordMention() string {
	if user == nil {
		return "<NULL>"
	}

	return fmt.Sprintf("<@%s>", user.ID)
}
