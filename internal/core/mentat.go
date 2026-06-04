package core

import (
	"fmt"

	"mentat/internal/config"
	"mentat/internal/history"
)

type Mentat struct {
	VaultDir string
	History  history.Store
	Decks    []*Deck
}

func NewMentat() (*Mentat, error) {
	vaultDir, err := config.Vault()
	if err != nil {
		return nil, fmt.Errorf("couldn't load vault dir: %w", err)
	}

	mentat := Mentat{
		Decks:    make([]*Deck, 0),
		VaultDir: vaultDir,
	}

	historyStore, err := history.NewFileStore(vaultDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create history store: %w", err)
	}
	mentat.History = historyStore

	// Load decks with review history
	decks, err := NewDecksFromDir(vaultDir, historyStore)
	if err != nil {
		return nil, fmt.Errorf("failed to load decks: %w", err)
	}
	mentat.Decks = decks

	return &mentat, nil
}

func (m *Mentat) Close() error {
	// Close history store when command completes
	if m.History != nil {
		if err := m.History.Close(); err != nil {
			return fmt.Errorf("failed to close history store: %w", err)
		}
	}
	return nil
}
