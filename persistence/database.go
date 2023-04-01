package persistence

type Database interface {
	GetDiscordMessage(id string) (string, error)
	SetDiscordMessage(id string, messageID string) error
}
