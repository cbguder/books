package cmd

import (
	"context"
	"github.com/cbguder/books/overdrive"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for media",
	Args:  cobra.ExactArgs(1),
	RunE:  search,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("library", "l", "", "Library code (e.g. sfpl)")
	searchCmd.MarkFlagRequired("library")
}

func search(cmd *cobra.Command, args []string) error {
	library, err := cmd.Flags().GetString("library")
	if err != nil {
		return err
	}

	client := overdrive.NewClient(cfg.Identity)
	resp, err := client.GetMedia(context.Background(), library, args[0])
	if err != nil {
		return err
	}

	t := newTableWriter()
	t.AppendHeader(table.Row{"Author", "Title", "Year", "Type", "Language", "Available", "Est. Wait"})

	for _, item := range resp.Items {
		t.AppendRow(table.Row{
			item.FirstCreatorName,
			item.Title,
			item.PublishDate.Year(),
			item.Type.Name,
			item.Languages[0].Name,
			item.IsAvailable,
			item.EstimatedWaitDays,
		})
	}

	t.Render()

	return nil
}
