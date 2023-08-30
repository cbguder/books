package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/cbguder/books/overdrive"
	"github.com/cbguder/books/repackage"
	"github.com/spf13/cobra"
)

var repackageCmd = &cobra.Command{
	Use:   "repackage <path>",
	Short: "repackage a book into an epub or MP3",
	Args:  cobra.ExactArgs(1),
	RunE:  repackageE,
}

func init() {
	rootCmd.AddCommand(repackageCmd)
}

func repackageE(_ *cobra.Command, args []string) error {
	srcDir := args[0]

	openbookPath := filepath.Join(srcDir, "_d", "openbook.json")

	openbook, err := overdrive.ReadOpenbook(openbookPath)
	if err != nil {
		return err
	}

	baseName := fmt.Sprintf("%s - %s", openbook.Creator[0].Name, openbook.Title.Main)

	if openbook.RenditionFormat == "ebook" {
		return repackage.Ebook(srcDir, baseName+".epub", openbook)
	} else if openbook.RenditionFormat == "audiobook" {
		return repackage.Audiobook(srcDir, baseName+".mp3", openbook)
	}

	return fmt.Errorf("unknown format: %s", openbook.RenditionFormat)
}
