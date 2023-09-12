package goodreads

import (
	"context"
	"strings"

	"github.com/cbguder/books/goodreads"
	"github.com/spf13/cobra"
)

var updateStatusCmd = &cobra.Command{
	Use:   "update-status <id> <page or percent> [body]",
	Short: "Update status for a book",
	Args:  cobra.RangeArgs(2, 3),
	RunE:  updateStatus,
}

func init() {
	GoodreadsCmd.AddCommand(updateStatusCmd)
}

func updateStatus(_ *cobra.Command, args []string) error {
	client := goodreads.NewClient()

	var page, percent, body string
	bookId := args[0]

	if strings.HasSuffix(args[1], "%") {
		percent = strings.TrimSuffix(args[1], "%")
	} else {
		page = args[1]
	}

	if len(args) > 2 {
		body = args[2]
	}

	return client.UpdateUserStatus(context.Background(), bookId, page, percent, body)
}
