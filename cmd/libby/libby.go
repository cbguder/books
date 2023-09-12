package libby

import (
	"fmt"

	"github.com/cbguder/books/config"
	"github.com/spf13/cobra"
)

var LibbyCmd = &cobra.Command{
	Use:     "libby",
	Aliases: []string{"lib"},
	Short:   "Libby commands",

	PersistentPreRunE: ensureAuth,
}

func ensureAuth(cmd *cobra.Command, _ []string) error {
	if cmd.Name() == "auth" {
		return nil
	}

	cfg := config.Get()
	if cfg.Identity == "" {
		return fmt.Errorf("Libby access token not found. Please authenticate first.")
	}

	return nil
}
