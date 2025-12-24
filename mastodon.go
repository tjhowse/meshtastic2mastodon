package main

import (
	"context"

	"github.com/mattn/go-mastodon"
)

type Mastodon struct {
	c *mastodon.Client
}

func NewMastodon(server, id, secret, access_token string) (*Mastodon, error) {
	m := &Mastodon{}
	m.c = mastodon.NewClient(&mastodon.Config{
		Server:       server,
		ClientID:     id,
		ClientSecret: secret,
		AccessToken:  access_token,
	})
	return m, nil
}

// Posts a status update
func (m *Mastodon) PostStatus(status string) error {
	_, err := m.c.PostStatus(context.Background(), &mastodon.Toot{
		Status: status,
	})
	return err
}
