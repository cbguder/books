package cmd

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cbguder/books/overdrive"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download <id>",
	Short: "Download a borrowed media item",
	Args:  cobra.ExactArgs(1),
	RunE:  download,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}

func download(_ *cobra.Command, args []string) error {
	mediaId := args[0]

	loan, err := findLoan(context.Background(), mediaId)
	if err != nil {
		return err
	}

	client := overdrive.NewClient(cfg.Identity)
	ctx := context.Background()

	fmt.Println("Opening book...")
	bookResp, err := client.OpenBook(ctx, loan.CardId, mediaId, loan.Type.Id)
	if err != nil {
		return err
	}

	rosters, err := client.GetRosters(ctx, bookResp.Urls.Rosters, bookResp.Message)
	if err != nil {
		return err
	}

	contentRoster, err := findContentRoster(rosters)
	if err != nil {
		return err
	}

	destFolder := fmt.Sprintf("%s - %s", loan.FirstCreatorName, loan.Title)
	fmt.Printf("Downloading %d files to \"%s\"...\n", len(contentRoster.Entries)+1, destFolder)

	err = os.MkdirAll(destFolder, 0755)
	if err != nil {
		return err
	}

	for _, entry := range contentRoster.Entries {
		err = downloadToFile(ctx, client, destFolder, entry.Url)
		if err != nil {
			return err
		}
	}

	err = downloadToFile(ctx, client, destFolder, bookResp.Urls.Openbook)
	if err != nil {
		return err
	}

	return nil
}

func findContentRoster(rosters []overdrive.Roster) (*overdrive.Roster, error) {
	for _, roster := range rosters {
		if roster.Group == "title-content" {
			return &roster, nil
		}
	}

	return nil, fmt.Errorf("failed to find content roster")
}

func filenameFromUrl(entryUrl string) (string, error) {
	u, err := url.Parse(entryUrl)
	if err != nil {
		return "", err
	}

	return filepath.Base(u.Path), nil
}

func downloadToFile(ctx context.Context, client *overdrive.Client, destFolder, srcUrl string) error {
	fname, err := filenameFromUrl(srcUrl)
	if err != nil {
		return err
	}

	fpath := filepath.Join(destFolder, fname)

	f, err := os.Create(fpath)
	if err != nil {
		return err
	}

	defer f.Close()

	resp, err := client.Download(ctx, srcUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func findLoan(ctx context.Context, mediaId string) (*overdrive.Loan, error) {
	syncResp, err := sync(ctx)
	if err != nil {
		return nil, err
	}

	for _, loan := range syncResp.Loans {
		if loan.Id == mediaId {
			return &loan, nil
		}
	}

	return nil, fmt.Errorf("no loan found for %s", mediaId)
}
