package libby

import (
	"context"

	"github.com/cbguder/books/cards"
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
	ctx := context.Background()

	media, err := client.GetMedia(ctx, card.Library.Key, mediaId)
	if err != nil {
		return err
	}

	return client.CreateLoan(ctx, card.Id, mediaId, media.Type.Id)
}

func getCard(libraryCode string) (*config.Card, error) {
	if libraryCode == "" {
		return cards.GetDefault()
	}

	return cards.GetByLibrary(libraryCode)
}
