package libby

import (
	"context"

	"github.com/cbguder/books/overdrive"
	"github.com/spf13/cobra"
)

var returnCmd = &cobra.Command{
	Use:   "return <id>",
	Short: "Return a loaned item",
	Args:  cobra.ExactArgs(1),
	RunE:  returnE,
}

func init() {
	LibbyCmd.AddCommand(returnCmd)
}

func returnE(_ *cobra.Command, args []string) error {
	mediaId := args[0]

	ctx := context.Background()

	loan, err := findLoan(ctx, mediaId)
	if err != nil {
		return err
	}

	client := overdrive.NewClient()
	return client.DeleteLoan(ctx, loan.CardId, mediaId)
}
