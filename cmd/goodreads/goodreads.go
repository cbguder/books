package goodreads

import (
	"context"
	"fmt"
	"time"

	"github.com/cbguder/books/config"
	"github.com/cbguder/books/goodreads"
	"github.com/spf13/cobra"
)

var GoodreadsCmd = &cobra.Command{
	Use:   "goodreads",
	Short: "Goodreads commands",

	PersistentPreRunE: ensureAuth,
}

func ensureAuth(cmd *cobra.Command, _ []string) error {
	if cmd.Name() == "auth" {
		return nil
	}

	cfg := config.Get()
	if cfg.Goodreads.AccessToken == "" || cfg.Goodreads.RefreshToken == "" {
		return fmt.Errorf("Goodreads access token not found. Please authenticate first.")
	}

	expiration := time.Unix(cfg.Goodreads.ExpiresAt, 0)
	if expiration.Before(time.Now()) {
		return goodreadsRefreshToken()
	}

	return nil
}

func goodreadsRefreshToken() error {
	fmt.Println("Refreshing Goodreads access token...")

	client := goodreads.NewClient()
	resp, err := client.Token(context.Background())
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(resp.ExpiresIn))

	cfg := config.Get()
	cfg.Goodreads.AccessToken = resp.AccessToken
	cfg.Goodreads.ExpiresAt = expiresAt.Unix()
	return cfg.Save()
}
