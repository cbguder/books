package cmd

import (
	"os"
	"path/filepath"

	"github.com/cbguder/books/cmd/goodreads"
	"github.com/cbguder/books/cmd/libby"
	"github.com/cbguder/books/config"
	"github.com/spf13/cobra"
)

var (
	Version string
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:     "books",
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(goodreads.GoodreadsCmd)
	rootCmd.AddCommand(libby.LibbyCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.books.yml)")
}

func initConfig() {
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cfgFile = filepath.Join(home, ".books.yml")
	}

	config.Load(cfgFile)
}
