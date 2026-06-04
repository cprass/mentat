package core

import (
	"testing"

	"mentat/internal/history"
)

func TestIsContentType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		cType    contentType
		expected bool
	}{
		{"question exact match", "q", typeQuestion, true},
		{"answer exact match", "a", typeAnswer, true},
		{"cloze exact match", "c", typeCloze, true},
		{"note exact match", "n", typeNote, true},
		{"question with space prefix", "q ", typeQuestion, true},
		{"answer with space prefix", "a ", typeAnswer, true},
		{"wrong type", "q", typeAnswer, false},
		{"empty string", "", typeQuestion, false},
		{"random text", "question", typeQuestion, false},
		{"frontmatter has no symbol", "---", typeFrontmatter, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isContentType(tt.input, tt.cType)
			if result != tt.expected {
				t.Errorf("isContentType(%q, %d) = %v, want %v", tt.input, tt.cType, result, tt.expected)
			}
		})
	}
}

func TestIsQuestion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"exact q", "q", true},
		{"q with space", "q ", true},
		{"q with text", "q What is this?", true},
		{"not a question", "a", false},
		{"random text", "question", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isQuestion(tt.input)
			if result != tt.expected {
				t.Errorf("isQuestion(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsAnswer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"exact a", "a", true},
		{"a with space", "a ", true},
		{"a with text", "a The answer", true},
		{"not an answer", "q", false},
		{"random text", "answer", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAnswer(tt.input)
			if result != tt.expected {
				t.Errorf("isAnswer(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewDeckFromFile_Dune(t *testing.T) {
	deck, err := NewDeckFromFile("../../test/dune.md", "../../test", []*history.ReviewEvent{})
	if err != nil {
		t.Fatalf("NewDeckFromFile failed: %v", err)
	}

	// Check deck name from heading
	if deck.Name != "Deck Title" {
		t.Errorf("Expected deck name 'Deck Title', got %q", deck.Name)
	}

	// Check number of cards parsed
	expectedCardCount := 4 // q/a pair, standalone q/a, and cloze
	if len(deck.Cards) != expectedCardCount {
		t.Errorf("Expected %d cards, got %d", expectedCardCount, len(deck.Cards))
	}

	// Check first card (inline q/a)
	if len(deck.Cards) > 3 {
		card1 := deck.Cards[3]
		expectedQ1 := "What house was Paul from?"
		expectedA1 := "Atreides"

		if card1.QuestionText != expectedQ1 {
			t.Errorf("Card 1 question: expected %q, got %q", expectedQ1, card1.QuestionText)
		}
		if card1.AnswerText != expectedA1 {
			t.Errorf("Card 1 answer: expected %q, got %q", expectedA1, card1.AnswerText)
		}
	}

	// Check second card (inline cloze)
	if len(deck.Cards) > 2 {
		card2 := deck.Cards[2]
		if card2.ClozeText == "" {
			t.Error("Card 2 cloze text is empty")
		}
	}

	// Check third card (multi-line q/a)
	if len(deck.Cards) > 1 {
		card3 := deck.Cards[1]

		// Question should contain markdown formatting
		if card3.QuestionText == "" {
			t.Error("Card 3 question is empty")
		}

		// Answer should contain the quote
		if card3.AnswerText == "" {
			t.Error("Card 3 answer is empty")
		}

		t.Logf("Card 3 Question:\n%s", card3.QuestionText)
		t.Logf("Card 3 Answer:\n%s", card3.AnswerText)
	}

	// Check 4th card (multi-line cloze)
	if len(deck.Cards) > 0 {
		card4 := deck.Cards[0]
		if card4.ClozeText == "" {
			t.Error("Card 4 cloze text is empty")
		}
	}

	// Debug: print all cards
	t.Logf("Total cards parsed: %d", len(deck.Cards))
	for i, card := range deck.Cards {
		t.Logf("Card %d - Type: %d", i+1, card.Type)
		t.Logf("  Question: %q", card.QuestionText)
		t.Logf("  Answer: %q", card.AnswerText)
	}
}

func TestNewDeckFromFile_Golang(t *testing.T) {
	deck, err := NewDeckFromFile("../../test/langs/golang.md", "../../test", []*history.ReviewEvent{})
	if err != nil {
		t.Fatalf("NewDeckFromFile failed: %v", err)
	}

	// Check deck name from heading
	if deck.Name != "Golang" {
		t.Errorf("Expected deck name 'Golang', got %q", deck.Name)
	}

	if deck.Location != "langs/golang.md" {
		t.Errorf("Expected deck location 'langs/golang.md', got %q", deck.Location)
	}

	if deck.Frontmatter.Created.IsZero() {
		t.Errorf("Expected frontmatter created date to be set, got zero value")
	}
}
