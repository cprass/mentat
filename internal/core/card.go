package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"mentat/internal/history"

	"github.com/open-spaced-repetition/go-fsrs/v3"
)

type CardType uint8

const (
	CardTypeQuestion CardType = iota
	CardTypeCloze
)

type Card struct {
	ID           string
	Type         CardType
	QuestionText string
	AnswerText   string
	ClozeText    string
	ClozeValues  []string
	ReviewEvents []*history.ReviewEvent
	Deck         *Deck
}

func (c *Card) LastReview() (*history.ReviewEvent, bool) {
	if len(c.ReviewEvents) == 0 {
		return nil, false
	}
	return c.ReviewEvents[len(c.ReviewEvents)-1], true
}

func (c *Card) Due() time.Time {
	rev, ok := c.LastReview()
	if !ok {
		return time.Now()
	}
	return rev.Card.Due
}

func (c *Card) Review(rating fsrs.Rating) (*history.ReviewEvent, error) {
	return c.Deck.ReviewCard(c, rating)
}

// Returns the card text for the given step and a boolean
// indicating whether there is a next step
func (c *Card) Reveal(step int) (string, bool) {
	if step < 0 {
		return "", false
	}

	switch c.Type {
	case CardTypeQuestion:
		if step == 0 {
			return c.QuestionText, true
		}
		if step == 1 {
			return c.AnswerText, false
		}
	case CardTypeCloze:
		hasNext := step < len(c.ClozeValues)
		if step == 0 {
			return c.ClozeText, hasNext
		}
		if step > len(c.ClozeValues) {
			return "", false
		}
		text := c.ClozeText
		for i := range step {
			text = strings.Replace(text, "[...]", c.ClozeValues[i], 1)
		}
		return text, hasNext
	}

	return "", false
}

func (c *Card) Title(step int) string {
	switch c.Type {
	case CardTypeCloze:
		return fmt.Sprintf("cloze %d/%d", step, len(c.ClozeValues))
	case CardTypeQuestion:
		if step == 0 {
			return "question"
		}
		return "answer"
	}
	return "card"
}

func createID(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:16])
}

func NewQuestionCard(deck *Deck, question, answer string) Card {
	return Card{
		ID:           createID(question),
		Type:         CardTypeQuestion,
		QuestionText: question,
		AnswerText:   answer,
		ReviewEvents: []*history.ReviewEvent{},
		Deck:         deck,
	}
}

func NewClozeCard(deck *Deck, cloze string) Card {
	card := Card{
		ID:           createID(cloze),
		Type:         CardTypeCloze,
		ReviewEvents: []*history.ReviewEvent{},
		Deck:         deck,
		ClozeValues:  []string{},
	}

	clozeRegexp := regexp.MustCompile(`\(\((.*?)\)\)`)
	matches := clozeRegexp.FindAllStringSubmatchIndex(cloze, -1)

	for _, match := range matches {
		clozeValue := cloze[match[2]:match[3]]
		card.ClozeValues = append(card.ClozeValues, clozeValue)
	}

	card.ClozeText = clozeRegexp.ReplaceAllString(cloze, "[...]")

	return card
}
