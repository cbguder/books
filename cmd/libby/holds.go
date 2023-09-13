package libby

import (
	"context"
	"fmt"

	"github.com/cbguder/books/cmd/out"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var holdsCmd = &cobra.Command{
	Use:   "holds",
	Short: "Show current holds",
	RunE:  holds,
}

func init() {
	LibbyCmd.AddCommand(holdsCmd)
}

func holds(_ *cobra.Command, _ []string) error {
	resp, err := sync(context.Background())
	if err != nil {
		return err
	}

	t := out.NewTableWriter()
	t.AppendHeader(table.Row{"ID", "Author", "Title", "Type", "Hold Placed", "Est. Wait"})

	for _, hold := range resp.Holds {
		t.AppendRow(table.Row{
			hold.Id,
			hold.FirstCreatorName,
			hold.Title,
			hold.Type.Name,
			hold.PlacedDate.Format("2006-01-02"),
			estimatedWaitStr(hold.EstimatedWaitDays),
		})
	}

	t.Render()

	return nil
}

func estimatedWaitStr(days int) string {
	if days == 0 {
		return "Unknown wait"
	}

	if days == 1 {
		return "1 day"
	}

	return fmt.Sprintf("%d days", days)
}
