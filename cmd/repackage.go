package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/cbguder/books/overdrive"
	"github.com/cbguder/books/repackage"
)

var repackageCmd = &cobra.Command{
	Use:   "repackage <path>",
	Short: "Repackage a book into an epub or MP3",
	Args:  cobra.ExactArgs(1),
	RunE:  repackageE,
}

func init() {
	rootCmd.AddCommand(repackageCmd)

	repackageCmd.Flags().Bool("fallback-cover", false, "use fallback cover image")
}

func repackageE(cmd *cobra.Command, args []string) error {
	srcDir := args[0]

	openbookPath := filepath.Join(srcDir, "_d", "openbook.json")

	openbook, err := overdrive.ReadOpenbook(openbookPath)
	if err != nil {
		return err
	}

	baseName := fmt.Sprintf("%s - %s", openbook.Creator[0].Name, openbook.Title.Main)

	if openbook.RenditionFormat == "ebook" {
		fallbackCover, _ := cmd.Flags().GetBool("fallback-cover")

		opts := repackage.EbookOptions{
			FallbackCover: fallbackCover,
		}

		return repackage.Ebook(srcDir, baseName+".epub", openbook, opts)
	} else if openbook.RenditionFormat == "audiobook" {
		return repackage.Audiobook(srcDir, baseName+".mp3", openbook)
	}

	return fmt.Errorf("unknown format: %s", openbook.RenditionFormat)
}
