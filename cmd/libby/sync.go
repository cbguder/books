package libby

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cbguder/books/config"
	"github.com/cbguder/books/overdrive"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Force sync",
	RunE:  syncE,
}

func init() {
	LibbyCmd.AddCommand(syncCmd)
}

func syncE(_ *cobra.Command, _ []string) error {
	_, err := sync(context.Background())
	return err
}

func sync(ctx context.Context) (*overdrive.SyncResponse, error) {
	client := overdrive.NewClient()

	resp, err := client.ChipSync(ctx)
	if err != nil {
		return nil, err
	}

	cfg := config.Get()
	cfg.Identity = resp.Identity
	cfg.Cards = mapCards(resp.Cards)
	err = cfg.Save()

	return resp, err
}

func mapCards(cards []overdrive.Card) []config.Card {
	cfgCards := make([]config.Card, len(cards))

	for i, card := range cards {
		cfgCards[i] = config.Card{
			Id:   card.CardId,
			Name: card.CardName,
			Library: config.Library{
				Name: card.Library.Name,
				Key:  card.AdvantageKey,
			},
		}
	}
	return cfgCards
}
