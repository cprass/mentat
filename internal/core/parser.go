package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"mentat/internal/history"

	"github.com/goccy/go-yaml"
)

type contentType = int

const (
	typeNone contentType = iota
	typeQuestion
	typeAnswer
	typeCloze
	typeNote
	typeFrontmatter
)

const (
	questionPrefix    = "q "
	answerPrefix      = "a "
	clozePrefix       = "c "
	notePrefix        = "n "
	frontmatterSymbol = "---"
	headingPrefix     = "# "
)

var contentTypeSymbolMap = map[contentType]string{
	typeQuestion: "q",
	typeAnswer:   "a",
	typeCloze:    "c",
	typeNote:     "n",
}

type parser struct {
	deck        *Deck
	currentCard *Card
	currentType contentType
	content     strings.Builder
	lineNum     int
	cardIDMap   map[string]*Card
}

func (p *parser) finishCurrentSection() error {
	text := strings.TrimSpace(p.content.String())
	if text == "" && p.currentType != typeNone {
		return nil
	}

	switch p.currentType {
	case typeQuestion:
		qCard := NewQuestionCard(p.deck, text, "")
		p.currentCard = &qCard
	case typeAnswer:
		p.currentCard.AnswerText = text
		p.commitCard()
	case typeCloze:
		cCard := NewClozeCard(p.deck, text)
		p.currentCard = &cCard
		p.commitCard()
	case typeNote:
	// TODO: handle notes
	case typeFrontmatter:
		var fm Frontmatter

		if err := yaml.Unmarshal([]byte(text), &fm); err != nil {
			return err
		}
		p.deck.Frontmatter = fm
	}

	p.content.Reset()
	return nil
}

func (p *parser) commitCard() {
	if p.currentCard != nil && (p.currentCard.QuestionText != "" || p.currentCard.ClozeText != "") {
		p.deck.Cards = append(p.deck.Cards, p.currentCard)
		p.cardIDMap[p.currentCard.ID] = p.currentCard
	}
}

func (p *parser) addLine(l string, trim bool) {
	if trim {
		switch {
		case isQuestion(l):
			l = strings.TrimPrefix(l, questionPrefix)
		case isAnswer(l):
			l = strings.TrimPrefix(l, answerPrefix)
		case isNote(l):
			l = strings.TrimPrefix(l, notePrefix)
		case isCloze(l):
			l = strings.TrimPrefix(l, clozePrefix)
		}
	}
	p.content.WriteByte('\n')
	p.content.WriteString(l)
}

func (p *parser) validateTransition(from, to contentType) error {
	if from == typeQuestion {
		if to != typeAnswer {
			return fmt.Errorf("question must be followed by answer, got %s", contentTypeSymbolMap[to])
		}
	}
	if from == typeAnswer && to == typeAnswer {
		return fmt.Errorf("answer can't be followed by answer")
	}
	return nil
}

func (p *parser) processLine(line string) error {
	p.lineNum++

	// Handle frontmatter
	if p.lineNum == 1 && line == frontmatterSymbol {
		p.currentType = typeFrontmatter
		return nil
	}
	if p.currentType == typeFrontmatter {
		if line == frontmatterSymbol {
			if err := p.finishCurrentSection(); err != nil {
				return err
			}
			p.currentType = typeNone
			return nil
		}
		p.addLine(line, false)
		return nil
	}

	// Extract deck name from heading
	if p.currentType == typeNone && strings.HasPrefix(line, headingPrefix) {
		if heading := strings.TrimSpace(strings.TrimPrefix(line, headingPrefix)); heading != "" {
			p.deck.Name = heading
		}
		return nil
	}

	// Detect section type changes
	var newType contentType
	switch {
	case isQuestion(line):
		newType = typeQuestion
	case isAnswer(line):
		newType = typeAnswer
	case isCloze(line):
		newType = typeCloze
	case isNote(line):
		newType = typeNote
	default:
		// Regular content line without new section marker
		p.addLine(line, false)
		return nil
	}

	// Validate transitions
	if err := p.validateTransition(p.currentType, newType); err != nil {
		return err
	}

	if err := p.finishCurrentSection(); err != nil {
		return err
	}
	p.currentType = newType

	// Handle inline content (e.g., "q What is this?")
	if len(line) > 2 {
		trimmed := strings.TrimSpace(line[2:]) // skip "q " or "a "
		if trimmed != "" {
			p.content.WriteString(trimmed + "\n")
		}
	}

	return nil
}

func isContentType(s string, c contentType) bool {
	if symbol, ok := contentTypeSymbolMap[c]; ok && (s == symbol || strings.HasPrefix(s, fmt.Sprintf("%s ", symbol))) {
		return true
	}
	return false
}

func isQuestion(line string) bool {
	return isContentType(line, typeQuestion)
}

func isAnswer(line string) bool {
	return isContentType(line, typeAnswer)
}

func isCloze(line string) bool {
	return isContentType(line, typeCloze)
}

func isNote(line string) bool {
	return isContentType(line, typeNote)
}

func NewDeckFromFile(path string, vaultDir string, reviews []*history.ReviewEvent) (*Deck, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	relPath, err := filepath.Rel(vaultDir, path)
	if err != nil {
		return nil, err
	}

	deck := &Deck{
		Cards:    []*Card{},
		Name:     strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Location: relPath,
	}

	parser := parser{
		deck:        deck,
		currentCard: &Card{},
		cardIDMap:   make(map[string]*Card),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		parser.processLine(text)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// adding an empty note to the end of the deck commits the last section
	parser.processLine("n")

	// add review data to cards
	for _, review := range reviews {
		card, exists := parser.cardIDMap[review.CardID]
		if !exists {
			continue
		}

		card.ReviewEvents = append(card.ReviewEvents, review)
	}

	// last card will be the first in the deck
	slices.Reverse(deck.Cards)

	// sort by
	// - cards with no reviews and latest card first
	// - cards with reviews sorted by smallest due-date first
	sort.SliceStable(deck.Cards, func(i, j int) bool {
		iHasReviews := len(deck.Cards[i].ReviewEvents) > 0
		jHasReviews := len(deck.Cards[j].ReviewEvents) > 0

		// Case 1: One has reviews, one doesn't - no reviews come first
		if iHasReviews != jHasReviews {
			return !iHasReviews // true if i has no reviews
		}

		// Case 2: Both have no reviews - keep order
		if !iHasReviews {
			return false
		}

		// Case 3: Both have reviews - sort by due date
		iLastRev, iOk := deck.Cards[i].LastReview()
		jLastRev, jOk := deck.Cards[j].LastReview()
		if !iOk || !jOk {
			return false
		}
		iDue := iLastRev.Card.Due
		jDue := jLastRev.Card.Due
		return iDue.Before(jDue)
	})

	return deck, nil
}
