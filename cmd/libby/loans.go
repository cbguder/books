package libby

import (
	"context"

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
	t.AppendHeader(table.Row{"ID", "Author", "Title", "Type", "Checkout Date", "Due Date"})

	for _, loan := range resp.Loans {
		t.AppendRow(table.Row{
			loan.Id,
			loan.FirstCreatorName,
			loan.Title,
			loan.Type.Name,
			loan.CheckoutDate.Format("2006-01-02"),
			loan.ExpireDate.Format("2006-01-02"),
		})
	}

	t.Render()

	return nil
}
