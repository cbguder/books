package libby

import (
	"context"
	"fmt"

	"github.com/cbguder/books/config"
	"github.com/cbguder/books/overdrive"
	"github.com/spf13/cobra"
)

var borrowCmd = &cobra.Command{
	Use:   "borrow <id>",
	Short: "Borrow a media item",
	Args:  cobra.ExactArgs(1),
	RunE:  borrow,
}

func init() {
	LibbyCmd.AddCommand(borrowCmd)

	borrowCmd.Flags().StringP("library", "l", "", "Library code (e.g. sfpl, defaults to first library)")
}

func borrow(cmd *cobra.Command, args []string) error {
	mediaId := args[0]

	libraryFlag, _ := cmd.Flags().GetString("library")
	card, err := getCard(libraryFlag)
	if err != nil {
		return err
	}

	client := overdrive.NewClient()
	return client.CreateLoan(context.Background(), card.Id, mediaId, "audiobook")
}

func getCard(libraryCode string) (*config.Card, error) {
	cfg := config.Get()

	if len(cfg.Cards) == 0 {
		return nil, fmt.Errorf("no library cards stored")
	}

	if libraryCode == "" {
		return &cfg.Cards[0], nil
	}

	for _, card := range cfg.Cards {
		if card.Library.Key == libraryCode {
			return &card, nil
		}
	}

	return nil, fmt.Errorf("no library card found for %s", libraryCode)
}
