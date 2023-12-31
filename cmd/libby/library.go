package libby

import (
	"context"

	"github.com/cbguder/books/cmd/out"
	"github.com/cbguder/books/libby"
	"github.com/cbguder/books/overdrive"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var libraryCmd = &cobra.Command{
	Use:   "library <query>",
	Short: "Search for a library",
	Args:  cobra.ExactArgs(1),
	RunE:  library,
}

func init() {
	LibbyCmd.AddCommand(libraryCmd)
}

func library(_ *cobra.Command, args []string) error {
	libbyClient := libby.NewClient()
	acResp, err := libbyClient.Autocomplete(context.Background(), args[0])
	if err != nil {
		return err
	}

	websiteIds := uniqueWebsiteIds(acResp)

	odClient := overdrive.NewClient()
	libResp, err := odClient.GetLibrariesByWebsiteId(context.Background(), websiteIds)
	if err != nil {
		return err
	}

	t := out.NewTableWriter()
	t.AppendHeader(table.Row{"ID", "Name"})

	for _, item := range libResp.Items {
		t.AppendRow(table.Row{
			item.Id,
			item.Name,
		})
	}

	t.Render()

	return nil
}

func uniqueWebsiteIds(resp *libby.AutocompleteResponse) []int {
	var websiteIds []int
	seenWebsiteIds := make(map[int]struct{})

	for _, branch := range resp.Branches {
		for _, system := range branch.Systems {
			if _, ok := seenWebsiteIds[system.WebsiteId]; ok {
				continue
			}

			websiteIds = append(websiteIds, system.WebsiteId)
			seenWebsiteIds[system.WebsiteId] = struct{}{}
		}
	}

	return websiteIds
}
