package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
)

func (m *Mastodon) Toot(message string) error {
	_, err := m.client.PostStatus(context.Background(), &mastodon.Toot{
		Status: message,
	})
	return err
}
