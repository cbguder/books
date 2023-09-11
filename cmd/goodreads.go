package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/cbguder/books/config"
	"github.com/cbguder/books/goodreads"
	"github.com/spf13/cobra"
)

var goodreadsCmd = &cobra.Command{
	Use:   "goodreads",
	Short: "Goodreads",

	PersistentPreRunE: goodreadsEnsureToken,
}

func init() {
	rootCmd.AddCommand(goodreadsCmd)
}

func goodreadsEnsureToken(cmd *cobra.Command, _ []string) error {
	if cmd.Name() == "auth" {
		return nil
	}

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

	client := goodreads.NewClient(cfg.Goodreads.AccessToken, cfg.Goodreads.RefreshToken)
	resp, err := client.Token(context.Background())
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(resp.ExpiresIn))

	cfg.Goodreads.AccessToken = resp.AccessToken
	cfg.Goodreads.ExpiresAt = expiresAt.Unix()

	return config.WriteConfig(cfgFile, cfg)
}
