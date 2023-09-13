package libby

import (
	"context"

	"github.com/cbguder/books/overdrive"
	"github.com/spf13/cobra"
)

var holdCmd = &cobra.Command{
	Use:   "hold <id>",
	Short: "Hold a media item",
	Args:  cobra.ExactArgs(1),
	RunE:  hold,
}

func init() {
	LibbyCmd.AddCommand(holdCmd)

	holdCmd.Flags().StringP("library", "l", "", "Library code (e.g. sfpl, defaults to first library)")
}

func hold(cmd *cobra.Command, args []string) error {
	mediaId := args[0]

	libraryFlag, _ := cmd.Flags().GetString("library")
	card, err := getCard(libraryFlag)
	if err != nil {
		return err
	}

	client := overdrive.NewClient()
	return client.CreateHold(context.Background(), card.Id, mediaId)
}
