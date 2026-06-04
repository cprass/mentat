package core

import (
	"cmp"
	"fmt"
	"slices"
	"sync"
	"time"

	"mentat/internal/history"

	"github.com/open-spaced-repetition/go-fsrs/v3"
	"golang.org/x/sync/errgroup"
)

type Deck struct {
	Name        string
	Location    string
	Cards       []*Card
	Frontmatter Frontmatter
}

type Frontmatter struct {
	Created time.Time `yaml:"created,omitempty"`
}

func NewDecksFromDir(vaultDir string, store history.Store) ([]*Deck, error) {
	reviewEvents, err := store.LoadAll()
	if err != nil {
		return nil, err
	}

	paths, err := LoadFiles(vaultDir, ".md")
	if err != nil {
		return nil, fmt.Errorf("failed to load deck files: %w", err)
	}

	var g errgroup.Group
	var mu sync.Mutex
	var decks []*Deck

	for _, path := range paths {
		g.Go(func() error {
			deck, err := NewDeckFromFile(path, vaultDir, reviewEvents)
			if err != nil {
				return err
			}

			mu.Lock()
			decks = append(decks, deck)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	slices.SortFunc(decks, func(a, b *Deck) int {
		return cmp.Compare(a.Location, b.Location)
	})

	return decks, nil
}

func (d *Deck) NextCard() (next *Card, hasNext bool) {
	if len(d.Cards) == 0 {
		return
	}
	next, hasNext = d.Cards[0], true
	return
}

func (d *Deck) ReviewCard(card *Card, rating fsrs.Rating) (*history.ReviewEvent, error) {

	p := fsrs.DefaultParam()
	f := fsrs.NewFSRS(p)

	var fsrsCard fsrs.Card
	lastRev, ok := card.LastReview()
	if !ok {
		fsrsCard = fsrs.NewCard()
	} else {
		fsrsCard = lastRev.Card
	}

	schedulingInfo := f.Next(fsrsCard, time.Now(), rating)
	reviewEvent := history.NewReviewEvent(card.ID, rating, schedulingInfo.Card)
	card.ReviewEvents = append(card.ReviewEvents, reviewEvent)

	// remove card from the original deck position
	cardIdx := slices.IndexFunc(d.Cards, func(c *Card) bool {
		return c.ID == card.ID
	})
	if cardIdx == -1 {
		return nil, fmt.Errorf("could not determine last card position in deck")
	}
	d.Cards = slices.Delete(d.Cards, cardIdx, cardIdx+1)

	// Find the first card idx that has a larger due date
	nextCardIdx := slices.IndexFunc(d.Cards, func(c *Card) bool {
		lastRev, ok := c.LastReview()
		if !ok {
			return false
		}
		return lastRev.Card.Due.After(schedulingInfo.Card.Due)
	})
	if nextCardIdx == -1 {
		// add card to the end
		d.Cards = append(d.Cards, card)
	} else {
		// add card before the next-card
		d.Cards = slices.Insert(d.Cards, nextCardIdx, card)
	}

	return reviewEvent, nil
}

func (d *Deck) String() string {
	str := fmt.Sprintf("\"%s\" has %d cards", d.Name, len(d.Cards))
	card, ok := d.NextCard()
	if ok {
		due := card.Due()
		minutesUntilDue := int(time.Until(due).Minutes())
		daysUntilDue := minutesUntilDue / 60 / 24

		if minutesUntilDue == 0 {
			str += " - next card due now"
		} else if daysUntilDue < 1 {
			str += " - next card due soon"
		} else {
			str += fmt.Sprintf(" - next card due in %d days", daysUntilDue)
		}
	}

	return str
}

func (d *Deck) NextDueCard() (card *Card, ok bool) {
	next, hasNext := d.NextCard()
	if !hasNext {
		return
	}
	due := next.Due()
	if due.After(time.Now()) {
		return
	}
	card, ok = next, true
	return
}

func (d *Deck) CardsDue() int {
	n := 0
	for _, card := range d.Cards {
		review, ok := card.LastReview()
		if !ok || review.Card.Due.Before(time.Now().Truncate(time.Hour*24).Add(time.Hour*24)) {
			n += 1
			continue
		}
		// Cards are ordered. Once we see the first one that is not due, the loop can stop
		break
	}
	return n
}

func (d *Deck) Stats() DeckStats {
	var (
		newCards, learning, relearning, review int
		lastReviewDate                         time.Time
	)

	for _, c := range d.Cards {
		rev, ok := c.LastReview()
		if !ok {
			newCards += 1
			continue
		}

		if lastReviewDate.Before(rev.Card.LastReview) {
			lastReviewDate = rev.Card.LastReview
		}

		switch rev.Card.State {
		case fsrs.New:
			newCards += 1
		case fsrs.Learning:
			learning += 1
		case fsrs.Relearning:
			relearning += 1
		case fsrs.Review:
			review += 1
		}
	}

	return DeckStats{
		Cards: CardStats{
			Count:      len(d.Cards),
			Due:        d.CardsDue(),
			New:        newCards,
			Learning:   learning,
			Relearning: relearning,
			Review:     review,
		},
		LastReview: lastReviewDate,
	}
}

func (d *Deck) Reviews() int {
	n := 0
	for _, c := range d.Cards {
		n += len(c.ReviewEvents)
	}
	return n
}

func CardsDue(d []*Deck) int {
	n := 0
	for _, deck := range d {
		n += deck.CardsDue()
	}
	return n
}

func FindNextCardInDecks(decks []*Deck) (*Card, bool) {
	var selectedCard *Card

	for _, deck := range decks {
		card, ok := deck.NextCard()
		if !ok {
			continue
		}
		if selectedCard == nil || selectedCard.Due().After(card.Due()) {
			selectedCard = card
			continue
		}
	}

	return selectedCard, selectedCard != nil
}
