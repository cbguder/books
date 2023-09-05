package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	RunE:  version,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func version(_ *cobra.Command, _ []string) error {
	fmt.Println("books version", Version)
	return nil
}
