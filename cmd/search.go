package cmd

import (
	"context"
	"fmt"

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
}

func search(cmd *cobra.Command, args []string) error {
	client := overdrive.NewClient(cfg.Identity)

	libraryFlag, err := cmd.Flags().GetString("library")
	if err != nil {
		return err
	}

	if libraryFlag != "" {
		return searchSingleLibrary(client, "", libraryFlag, args[0])
	}

	if len(cfg.Cards) == 0 {
		return fmt.Errorf("please authenticate or specify a library")
	}

	for _, card := range cfg.Cards {
		err = searchSingleLibrary(client, card.LibraryName, card.LibraryKey, args[0])
		if err != nil {
			return err
		}
	}

	return nil
}

func searchSingleLibrary(client *overdrive.Client, title, libraryKey, query string) error {
	resp, err := client.GetMedia(context.Background(), libraryKey, query)
	if err != nil {
		return err
	}

	t := newTableWriter()

	if title != "" {
		t.SetTitle(title)
	}

	t.AppendHeader(table.Row{"ID", "Author", "Title", "Year", "Type", "Language", "Available", "Est. Wait"})

	for _, item := range resp.Items {
		t.AppendRow(table.Row{
			item.Id,
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
