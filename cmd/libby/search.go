package libby

import (
	"context"
	"fmt"

	"github.com/cbguder/books/cmd/out"
	"github.com/cbguder/books/config"
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
	LibbyCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("library", "l", "", "Library code (e.g. sfpl)")
	searchCmd.Flags().StringP("format", "f", "", "Media format (audiobook or ebook)")
}

func search(cmd *cobra.Command, args []string) error {
	client := overdrive.NewClient()

	libraryFlag, _ := cmd.Flags().GetString("library")

	formatFlag, _ := cmd.Flags().GetString("format")
	format := formatFromFlag(formatFlag)

	if libraryFlag != "" {
		return searchSingleLibrary(client, "", libraryFlag, args[0], format)
	}

	cfg := config.Get()
	if len(cfg.Cards) == 1 {
		card := cfg.Cards[0]
		return searchSingleLibrary(client, "", card.Library.Key, args[0], format)
	}

	if len(cfg.Cards) == 0 {
		return fmt.Errorf("please authenticate or specify a library")
	}

	for _, card := range cfg.Cards {
		err := searchSingleLibrary(client, card.Library.Name, card.Library.Key, args[0], format)
		if err != nil {
			return err
		}
	}

	return nil
}

func formatFromFlag(formatFlag string) overdrive.MediaFormat {
	if formatFlag == "audiobook" {
		return overdrive.MediaFormatAudiobook
	} else if formatFlag == "ebook" {
		return overdrive.MediaFormatEbook
	}

	return overdrive.MediaFormatAny
}

func searchSingleLibrary(client *overdrive.Client, title, libraryKey, query string, format overdrive.MediaFormat) error {
	resp, err := client.GetMedia(context.Background(), libraryKey, query, format)
	if err != nil {
		return err
	}

	t := out.NewTableWriter()

	if title != "" {
		t.SetTitle(title)
	}

	t.AppendHeader(table.Row{"ID", "Author", "Title", "Year", "Type", "Language", "Availability"})

	for _, item := range resp.Items {
		t.AppendRow(table.Row{
			item.Id,
			item.FirstCreatorName,
			item.Title,
			item.PublishDate.Year(),
			item.Type.Name,
			item.Languages[0].Name,
			availabilityStr(item.IsAvailable, item.EstimatedWaitDays),
		})
	}

	t.Render()

	return nil
}

func availabilityStr(available bool, estWait int) string {
	if available {
		return "Available"
	}

	if estWait == 0 {
		return "Unknown wait"
	}

	return fmt.Sprintf("%d day wait", estWait)
}
