package review

import (
	"fmt"
	"time"

	"mentat/internal/core"
	"mentat/internal/history"

	"github.com/open-spaced-repetition/go-fsrs/v3"
	"github.com/rivo/tview"
)

const oneDay = time.Hour * 24

// View holds the tview primitives for a review screen.
type View struct {
	Flex   *tview.Flex
	Text   *tview.TextView
	Footer *tview.TextView
}

func NewView() *View {
	text := tview.NewTextView()
	text.SetBorderPadding(1, 1, 2, 2).SetBorder(true)

	footer := tview.NewTextView()

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(text, 0, 1, true)
	flex.AddItem(footer, 1, 0, false)

	return &View{Flex: flex, Text: text, Footer: footer}
}

// CardReview manages the review state for a sequence of cards.
type CardReview struct {
	CardsProcessed int
	RevealIndex    int
	HasReveal      bool
	Card           *core.Card
	history        history.Store
	decks          []*core.Deck
	view           *View
}

// NewCardReview creates a review session and renders the first card.
// Returns nil if no cards are due.
func NewCardReview(hist history.Store, decks []*core.Deck, view *View) *CardReview {
	card := LoadNextDueCard(decks)
	if card == nil {
		return nil
	}

	cr := &CardReview{
		view:           view,
		decks:          decks,
		Card:           card,
		CardsProcessed: 1,
		HasReveal:      true,
		history:        hist,
	}
	cr.UpdateView()
	return cr
}

// UpdateView refreshes the text, title, and footer based on current state.
func (r *CardReview) UpdateView() {
	text, hasNext := r.Card.Reveal(r.RevealIndex)
	r.HasReveal = hasNext
	r.view.Text.SetText(text)

	title := r.Card.Title(r.RevealIndex)
	cardsRemaining := core.CardsDue(r.decks)
	r.view.Text.SetTitle(fmt.Sprintf(
		"deck %s - %s - %d cards left",
		r.Card.Deck.Location, title, cardsRemaining,
	))

	if hasNext {
		r.view.Footer.SetText("space:reveal - q:quit")
	} else {
		r.view.Footer.SetText("1:again - 2:hard - 3:good - 4:easy - q:quit")
	}
}

// Reveal advances to the next reveal step and updates the view.
func (r *CardReview) Reveal() {
	r.RevealIndex++
	r.UpdateView()
}

// Next submits a rating, persists it, and loads the next due card.
// Returns true if a next card was loaded, false if the session is done.
func (r *CardReview) Next(rating fsrs.Rating) bool {
	re, err := r.Card.Review(rating)
	if err != nil {
		return false
	}
	r.history.Append(re)

	nextCard := LoadNextDueCard(r.decks)
	if nextCard == nil {
		return false
	}

	r.CardsProcessed++
	r.Card = nextCard
	r.RevealIndex = 0
	r.UpdateView()
	return true
}

// LoadNextDueCard finds the next card due today across all decks.
// Returns nil if no cards are due before end of day.
func LoadNextDueCard(decks []*core.Deck) *core.Card {
	nextCard, _ := core.FindNextCardInDecks(decks)
	endOfDay := time.Now().Truncate(oneDay).Add(oneDay)

	if nextCard == nil || nextCard.Due().After(endOfDay) {
		return nil
	}
	return nextCard
}
