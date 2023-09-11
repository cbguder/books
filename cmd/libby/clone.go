package libby

import (
	"context"
	"fmt"
	"regexp"

	"github.com/cbguder/books/config"
	"github.com/cbguder/books/overdrive"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Libby",
	RunE:  auth,
}

func init() {
	LibbyCmd.AddCommand(cloneCmd)
}

func auth(_ *cobra.Command, _ []string) error {
	fmt.Println("Go to Menu > Settings > Copy To Another Device. You will see a setup code. Enter it below.")
	fmt.Print("Setup code: ")

	var code string
	_, err := fmt.Scanln(&code)
	if err != nil {
		return err
	}

	match, _ := regexp.MatchString(`^\d{8}$`, code)
	if !match {
		return fmt.Errorf("setup code must be 8 digits")
	}

	client := overdrive.NewClient()
	ctx := context.Background()

	_, err = client.Chip(ctx)
	if err != nil {
		return err
	}

	_, err = client.ChipClone(ctx, code)
	if err != nil {
		return err
	}

	resp, err := client.Chip(ctx)
	if err != nil {
		return err
	}

	cfg := config.Get()
	cfg.Identity = resp.Identity

	_, err = sync(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Clone successful")

	return nil
}
