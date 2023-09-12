package goodreads

import (
	"context"

	"github.com/cbguder/books/cmd/out"
	"github.com/cbguder/books/goodreads"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var shelfCmd = &cobra.Command{
	Use:   "shelf <name>",
	Short: "List books on shelf",
	Args:  cobra.ExactArgs(1),
	RunE:  shelf,
}

func init() {
	GoodreadsCmd.AddCommand(shelfCmd)
}

func shelf(_ *cobra.Command, args []string) error {
	client := goodreads.NewClient()

	resp, err := client.GetReviews(context.Background(), args[0])
	if err != nil {
		return err
	}

	t := out.NewTableWriter()
	t.AppendHeader(table.Row{"Author", "Title", "Year"})

	for _, review := range resp.Reviews {
		t.AppendRow(table.Row{
			review.Book.Authors[0].Name,
			review.Book.Title,
			review.Book.PublicationYear,
		})
	}

	t.Render()

	return nil
}
