package cmd

import (
	"context"

	"github.com/cbguder/books/goodreads"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var shelvesCmd = &cobra.Command{
	Use:   "shelves",
	Short: "List shelves",
	RunE:  shelvesE,
}

func init() {
	goodreadsCmd.AddCommand(shelvesCmd)
}

func shelvesE(_ *cobra.Command, _ []string) error {
	client := goodreads.NewClient(cfg.Goodreads.AccessToken, cfg.Goodreads.RefreshToken)

	resp, err := client.GetCurrentUserShelves(context.Background())
	if err != nil {
		return err
	}

	t := newTableWriter()
	t.AppendHeader(table.Row{"Name", "Book Count", "Exclusive?"})
	for _, shelf := range resp.CurrentUser.UserShelves {
		t.AppendRow(table.Row{shelf.Name, shelf.BookCount, shelf.ExclusiveFlag})
	}

	t.Render()

	return nil
}
