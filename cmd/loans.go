package cmd

import (
	"context"
	"os"

	"github.com/cbguder/books/overdrive"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var loansCmd = &cobra.Command{
	Use:   "loans",
	Short: "Show current loans",
	RunE:  loans,
}

func init() {
	rootCmd.AddCommand(loansCmd)
}

func loans(_ *cobra.Command, _ []string) error {
	client := overdrive.NewClient()

	resp, err := client.ChipSync(context.Background())
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	style := table.StyleRounded
	style.Format.Header = text.FormatDefault
	t.SetStyle(style)

	t.AppendHeader(table.Row{"Author", "Title", "Type", "Checkout Date", "Due Date"})

	for _, loan := range resp.Loans {
		t.AppendRow(table.Row{
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
