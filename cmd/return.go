package cmd

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
	rootCmd.AddCommand(returnCmd)
}

func returnE(_ *cobra.Command, args []string) error {
	mediaId := args[0]

	loan, err := findLoan(context.Background(), mediaId)
	if err != nil {
		return err
	}

	client := overdrive.NewClient(cfg.Identity)
	ctx := context.Background()

	return client.DeleteLoan(ctx, loan.CardId, mediaId)
}
