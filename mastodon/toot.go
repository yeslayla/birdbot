package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
)

// Toot publishes a toot on Mastodon
func (m *Mastodon) Toot(message string) error {
	_, err := m.client.PostStatus(context.Background(), &mastodon.Toot{
		Status: message,
	})
	return err
}
