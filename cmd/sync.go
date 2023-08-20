package cmd

import (
	"context"
	"github.com/cbguder/books/config"
	"github.com/cbguder/books/overdrive"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync data manually",
	RunE:  syncE,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func syncE(_ *cobra.Command, _ []string) error {
	_, err := sync(context.Background())
	return err
}

func sync(ctx context.Context) (*overdrive.SyncResponse, error) {
	client := overdrive.NewClient(cfg.Identity)

	resp, err := client.ChipSync(ctx)
	if err != nil {
		return nil, err
	}

	cfg.Identity = resp.Identity
	cfg.Cards = mapCards(resp.Cards)
	err = config.WriteConfig(cfgFile, cfg)

	return resp, err
}

func mapCards(cards []overdrive.Card) []config.Card {
	cfgCards := make([]config.Card, len(cards))

	for i, card := range cards {
		cfgCards[i] = config.Card{
			Id:          card.CardId,
			Name:        card.CardName,
			LibraryName: card.Library.Name,
			LibraryKey:  card.AdvantageKey,
		}
	}
	return cfgCards
}
