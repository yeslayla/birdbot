package mastodon

import (
	"context"
	"log"

	"github.com/mattn/go-mastodon"
)

type Mastodon struct {
	client *mastodon.Client
}

func NewMastodon(server string, clientID string, clientSecret string, username string, password string) *Mastodon {
	m := &Mastodon{}

	m.client = mastodon.NewClient(&mastodon.Config{
		Server:       server,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	err := m.client.Authenticate(context.Background(), username, password)
	if err != nil {
		log.Print("Failed to configure Mastodon: ", err)
		return nil
	}

	return m
}
