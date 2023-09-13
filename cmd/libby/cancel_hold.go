package libby

import (
	"context"

	"github.com/cbguder/books/overdrive"
	"github.com/spf13/cobra"
)

var cancelHoldCmd = &cobra.Command{
	Use:   "cancel-hold <id>",
	Short: "Cancel hold on a media item",
	Args:  cobra.ExactArgs(1),
	RunE:  cancelHold,
}

func init() {
	LibbyCmd.AddCommand(cancelHoldCmd)

	cancelHoldCmd.Flags().StringP("library", "l", "", "Library code (e.g. sfpl, defaults to first library)")
}

func cancelHold(cmd *cobra.Command, args []string) error {
	mediaId := args[0]

	libraryFlag, _ := cmd.Flags().GetString("library")
	card, err := getCard(libraryFlag)
	if err != nil {
		return err
	}

	client := overdrive.NewClient()
	return client.DeleteHold(context.Background(), card.Id, mediaId)
}
