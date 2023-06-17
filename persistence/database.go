package persistence

// Database is an interface used to wrap persistant data
type Database interface {
	GetDiscordMessage(id string) (string, error)
	SetDiscordMessage(id string, messageID string) error

	GetDiscordWebhook(id string) (*DBDiscordWebhook, error)
	SetDiscordWebhook(id string, data *DBDiscordWebhook) error
}

type DBDiscordWebhook struct {
	ID    string
	Token string
}
