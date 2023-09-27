package libby

import (
	"context"

	"github.com/cbguder/books/cards"
	"github.com/cbguder/books/cmd/out"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var loansCmd = &cobra.Command{
	Use:   "loans",
	Short: "Show current loans",
	RunE:  loans,
}

func init() {
	LibbyCmd.AddCommand(loansCmd)
}

func loans(_ *cobra.Command, _ []string) error {
	resp, err := sync(context.Background())
	if err != nil {
		return err
	}

	t := out.NewTableWriter()
	t.AppendHeader(table.Row{"ID", "Author", "Title", "Type", "Library", "Checkout Date", "Due Date"})

	for _, loan := range resp.Loans {
		card, err := cards.GetById(loan.CardId)
		if err != nil {
			return err
		}

		t.AppendRow(table.Row{
			loan.Id,
			loan.FirstCreatorName,
			loan.Title,
			loan.Type.Name,
			card.Library.Key,
			loan.CheckoutDate.Format("2006-01-02"),
			loan.ExpireDate.Format("2006-01-02"),
		})
	}

	t.Render()

	return nil
}
