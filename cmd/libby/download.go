package libby

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cbguder/books/overdrive"
)

var downloadCmd = &cobra.Command{
	Use:   "download <id>",
	Short: "Download a borrowed media item",
	Args:  cobra.ExactArgs(1),
	RunE:  download,
}

func init() {
	LibbyCmd.AddCommand(downloadCmd)
}

func download(_ *cobra.Command, args []string) error {
	mediaId := args[0]

	loan, err := findLoan(context.Background(), mediaId)
	if err != nil {
		return err
	}

	client := overdrive.NewClient()
	ctx := context.Background()

	bookResp, err := client.OpenLoan(ctx, loan.CardId, mediaId, loan.Type.Id)
	if err != nil {
		return err
	}

	destFolder := fmt.Sprintf("%s - %s", loan.FirstCreatorName, loan.Title)

	dl := overdrive.NewDownloader()
	return dl.Download(ctx, bookResp, destFolder)
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
