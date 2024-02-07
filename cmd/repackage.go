package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/cbguder/books/epubcheck"
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
	repackageCmd.Flags().Bool("check", false, "validate the output after repackaging")
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
		return repackageEbook(cmd, srcDir, baseName, openbook)
	} else if openbook.RenditionFormat == "audiobook" {
		return repackage.Audiobook(srcDir, baseName+".mp3", openbook)
	}

	return fmt.Errorf("unknown format: %s", openbook.RenditionFormat)
}

func repackageEbook(cmd *cobra.Command, srcDir, baseName string, openbook *overdrive.Openbook) error {
	fallbackCover, _ := cmd.Flags().GetBool("fallback-cover")

	opts := repackage.EbookOptions{
		FallbackCover: fallbackCover,
	}

	err := repackage.Ebook(srcDir, baseName+".epub", openbook, opts)
	if err != nil {
		return err
	}

	check, _ := cmd.Flags().GetBool("check")
	if check {
		return checkEbook(baseName + ".epub")
	}

	return nil
}

func checkEbook(path string) error {
	result, err := epubcheck.Check(path)
	if err != nil {
		return err
	}

	if result.HasErrors() {
		for _, msg := range result.Messages {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	return nil
}
