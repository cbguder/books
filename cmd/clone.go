package cmd

import (
	"context"
	"fmt"
	"regexp"

	"github.com/cbguder/books/overdrive"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone Libby account from another device",
	RunE:  clone,
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

func clone(_ *cobra.Command, _ []string) error {
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

	_, err = client.Chip(context.Background())
	if err != nil {
		return err
	}

	_, err = client.ChipClone(context.Background(), code)
	if err != nil {
		return err
	}

	_, err = client.Chip(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Clone successful")

	return nil
}
