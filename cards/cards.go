package cards

import (
	"fmt"

	"github.com/cbguder/books/config"
)

func GetDefault() (*config.Card, error) {
	cfg := config.Get()

	if len(cfg.Cards) == 0 {
		return nil, fmt.Errorf("no library cards stored")
	}

	return &cfg.Cards[0], nil
}

func GetByLibrary(libraryCode string) (*config.Card, error) {
	cfg := config.Get()

	if len(cfg.Cards) == 0 {
		return nil, fmt.Errorf("no library cards stored")
	}

	for _, card := range cfg.Cards {
		if card.Library.Key == libraryCode {
			return &card, nil
		}
	}

	return nil, fmt.Errorf("no library card found for %s", libraryCode)
}

func GetById(cardId string) (*config.Card, error) {
	cfg := config.Get()

	if len(cfg.Cards) == 0 {
		return nil, fmt.Errorf("no library cards stored")
	}

	for _, card := range cfg.Cards {
		if card.Id == cardId {
			return &card, nil
		}
	}

	return nil, fmt.Errorf("no library card found with id %s", cardId)
}
