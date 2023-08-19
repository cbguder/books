package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "books",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.books.yml)")
}

func initConfig() {
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cfgFile = filepath.Join(home, ".books.yml")
	}

	viper.SetConfigFile(cfgFile)
	viper.ReadInConfig()
}
