package cmd

import (
	"context"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var holdsCmd = &cobra.Command{
	Use:   "holds",
	Short: "Show current holds",
	RunE:  holds,
}

func init() {
	rootCmd.AddCommand(holdsCmd)
}

func holds(_ *cobra.Command, _ []string) error {
	resp, err := sync(context.Background())
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	style := table.StyleRounded
	style.Format.Header = text.FormatDefault
	t.SetStyle(style)

	t.AppendHeader(table.Row{"Author", "Title", "Type", "Hold Placed", "Est. Wait"})

	for _, hold := range resp.Holds {
		t.AppendRow(table.Row{
			hold.FirstCreatorName,
			hold.Title,
			hold.Type.Name,
			hold.PlacedDate.Format("2006-01-02"),
			hold.EstimatedWaitDays,
		})
	}

	t.Render()

	return nil
}
