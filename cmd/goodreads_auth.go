package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cbguder/books/config"
	"github.com/cbguder/books/goodreads"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var goodreadsAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate",
	RunE:  goodreadsAuth,
}

func init() {
	goodreadsCmd.AddCommand(goodreadsAuthCmd)
}

func goodreadsAuth(_ *cobra.Command, _ []string) error {
	client := goodreads.NewClient("", "")

	var email string

	fmt.Print("Email: ")
	_, err := fmt.Scanln(&email)
	if err != nil {
		return err
	}

	fmt.Print("Password: ")
	pwdBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Println("")

	password := string(pwdBytes)

	fmt.Println("Getting Goodreads access token...")

	ctx := context.Background()
	registerResp, err := client.Register(ctx, email, password)
	if err != nil {
		return err
	}

	expSecs, err := strconv.Atoi(registerResp.Response.Success.Tokens.Bearer.ExpiresIn)
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(expSecs))

	currentUserResp, err := client.GetCurrentUserData(ctx)
	if err != nil {
		return err
	}

	cfg.Goodreads.AccessToken = registerResp.Response.Success.Tokens.Bearer.AccessToken
	cfg.Goodreads.RefreshToken = registerResp.Response.Success.Tokens.Bearer.RefreshToken
	cfg.Goodreads.ExpiresAt = expiresAt.Unix()
	cfg.Goodreads.UserId = currentUserResp.CurrentUser.User.ID

	return config.WriteConfig(cfgFile, cfg)
}
