package core

import (
	"testing"
)

// --- Constructors ---

func TestNewQuestionCard(t *testing.T) {
	deck := &Deck{Name: "test"}
	card := NewQuestionCard(deck, "What is Go?", "A language.")

	if card.Type != CardTypeQuestion {
		t.Errorf("Type = %d, want CardTypeQuestion", card.Type)
	}
	if card.QuestionText != "What is Go?" {
		t.Errorf("QuestionText = %q, want %q", card.QuestionText, "What is Go?")
	}
	if card.AnswerText != "A language." {
		t.Errorf("AnswerText = %q, want %q", card.AnswerText, "A language.")
	}
	if card.Deck != deck {
		t.Error("Deck reference not set")
	}
}

func TestNewClozeCard(t *testing.T) {
	deck := &Deck{Name: "test"}

	tests := []struct {
		name           string
		input          string
		expectedText   string
		expectedValues []string
	}{
		{
			name:           "single cloze",
			input:          "((water)) is a very rare element on Dune.",
			expectedText:   "[...] is a very rare element on Dune.",
			expectedValues: []string{"water"},
		},
		{
			name:           "multiple clozes",
			input:          "((Fear)) is the ((mind))-killer.",
			expectedText:   "[...] is the [...]-killer.",
			expectedValues: []string{"Fear", "mind"},
		},
		{
			name:           "no clozes",
			input:          "Just a plain sentence.",
			expectedText:   "Just a plain sentence.",
			expectedValues: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := NewClozeCard(deck, tt.input)

			if card.Type != CardTypeCloze {
				t.Errorf("expected CardTypeCloze, got %d", card.Type)
			}

			if card.ClozeText != tt.expectedText {
				t.Errorf("ClozeText = %q, want %q", card.ClozeText, tt.expectedText)
			}

			if len(card.ClozeValues) != len(tt.expectedValues) {
				t.Fatalf("ClozeValues length = %d, want %d", len(card.ClozeValues), len(tt.expectedValues))
			}

			for i, v := range tt.expectedValues {
				if card.ClozeValues[i] != v {
					t.Errorf("ClozeValues[%d] = %q, want %q", i, card.ClozeValues[i], v)
				}
			}
		})
	}
}

// --- Card IDs ---

func TestCardID_Deterministic(t *testing.T) {
	deck := &Deck{Name: "test"}
	card1 := NewQuestionCard(deck, "same question", "a")
	card2 := NewQuestionCard(deck, "same question", "b")

	if card1.ID != card2.ID {
		t.Errorf("same question should produce same ID: %q != %q", card1.ID, card2.ID)
	}
}

func TestCardID_Unique(t *testing.T) {
	deck := &Deck{Name: "test"}
	card1 := NewQuestionCard(deck, "question one", "a")
	card2 := NewQuestionCard(deck, "question two", "a")

	if card1.ID == card2.ID {
		t.Errorf("different questions should produce different IDs: both %q", card1.ID)
	}
}

func TestClozeCardID_UsesRawInput(t *testing.T) {
	deck := &Deck{Name: "test"}
	card1 := NewClozeCard(deck, "((Fear)) is the mind-killer.")
	card2 := NewClozeCard(deck, "((Fear)) is the mind-killer.")

	if card1.ID != card2.ID {
		t.Errorf("same cloze input should produce same ID: %q != %q", card1.ID, card2.ID)
	}
}

// --- Reveal ---

func TestReveal_QuestionCard(t *testing.T) {
	deck := &Deck{Name: "test"}
	card := NewQuestionCard(deck, "What is Go?", "A programming language.")

	text, hasNext := card.Reveal(0)
	if text != "What is Go?" || !hasNext {
		t.Errorf("Reveal(0) = (%q, %v), want (%q, true)", text, hasNext, "What is Go?")
	}

	text, hasNext = card.Reveal(1)
	if text != "A programming language." || hasNext {
		t.Errorf("Reveal(1) = (%q, %v), want (%q, false)", text, hasNext, "A programming language.")
	}

	text, hasNext = card.Reveal(2)
	if text != "" || hasNext {
		t.Errorf("Reveal(2) = (%q, %v), want (%q, false)", text, hasNext, "")
	}
}

func TestReveal_ClozeCard(t *testing.T) {
	deck := &Deck{Name: "test"}
	card := NewClozeCard(deck, "((Fear)) is the ((mind))-killer.")

	text, hasNext := card.Reveal(0)
	if text != "[...] is the [...]-killer." || !hasNext {
		t.Errorf("Reveal(0) = (%q, %v), want (%q, true)", text, hasNext, "[...] is the [...]-killer.")
	}

	text, hasNext = card.Reveal(1)
	if text != "Fear is the [...]-killer." || !hasNext {
		t.Errorf("Reveal(1) = (%q, %v), want (%q, true)", text, hasNext, "Fear is the [...]-killer.")
	}

	text, hasNext = card.Reveal(2)
	if text != "Fear is the mind-killer." || hasNext {
		t.Errorf("Reveal(2) = (%q, %v), want (%q, false)", text, hasNext, "Fear is the mind-killer.")
	}

	text, hasNext = card.Reveal(3)
	if text != "" || hasNext {
		t.Errorf("Reveal(3) = (%q, %v), want (%q, false)", text, hasNext, "")
	}
}

func TestReveal_SingleCloze(t *testing.T) {
	deck := &Deck{Name: "test"}
	card := NewClozeCard(deck, "((water)) is rare on Dune.")

	text, hasNext := card.Reveal(0)
	if text != "[...] is rare on Dune." || !hasNext {
		t.Errorf("Reveal(0) = (%q, %v), want (%q, true)", text, hasNext, "[...] is rare on Dune.")
	}

	text, hasNext = card.Reveal(1)
	if text != "water is rare on Dune." || hasNext {
		t.Errorf("Reveal(1) = (%q, %v), want (%q, false)", text, hasNext, "water is rare on Dune.")
	}
}

func TestReveal_NegativeStep(t *testing.T) {
	deck := &Deck{Name: "test"}

	qCard := NewQuestionCard(deck, "q", "a")
	text, hasNext := qCard.Reveal(-1)
	if text != "" || hasNext {
		t.Errorf("Question Reveal(-1) = (%q, %v), want (%q, false)", text, hasNext, "")
	}

	cCard := NewClozeCard(deck, "((x)) y")
	text, hasNext = cCard.Reveal(-1)
	if text != "" || hasNext {
		t.Errorf("Cloze Reveal(-1) = (%q, %v), want (%q, false)", text, hasNext, "")
	}
}

func TestReveal_ClozeNoClozeValues(t *testing.T) {
	deck := &Deck{Name: "test"}
	card := NewClozeCard(deck, "No cloze markers here.")

	text, hasNext := card.Reveal(0)
	if text != "No cloze markers here." || hasNext {
		t.Errorf("Reveal(0) = (%q, %v), want (%q, false)", text, hasNext, "No cloze markers here.")
	}
}

func TestReveal_ClozeEmptyValue(t *testing.T) {
	deck := &Deck{Name: "test"}
	card := NewClozeCard(deck, "A (()) B")

	if len(card.ClozeValues) != 1 || card.ClozeValues[0] != "" {
		t.Fatalf("ClozeValues = %v, want [\"\"]", card.ClozeValues)
	}

	text, _ := card.Reveal(0)
	if text != "A [...] B" {
		t.Errorf("Reveal(0) = %q, want %q", text, "A [...] B")
	}

	text, hasNext := card.Reveal(1)
	if text != "A  B" || hasNext {
		t.Errorf("Reveal(1) = (%q, %v), want (%q, false)", text, hasNext, "A  B")
	}
}

// --- Title ---

func TestTitle(t *testing.T) {
	deck := &Deck{Name: "test"}

	qCard := NewQuestionCard(deck, "q", "a")
	if title := qCard.Title(0); title != "question" {
		t.Errorf("Question Title(0) = %q, want %q", title, "question")
	}
	if title := qCard.Title(1); title != "answer" {
		t.Errorf("Question Title(1) = %q, want %q", title, "answer")
	}

	cCard := NewClozeCard(deck, "((Fear)) is the ((mind))-killer.")
	if title := cCard.Title(0); title != "cloze 0/2" {
		t.Errorf("Cloze Title(0) = %q, want %q", title, "cloze 0/2")
	}
	if title := cCard.Title(1); title != "cloze 1/2" {
		t.Errorf("Cloze Title(1) = %q, want %q", title, "cloze 1/2")
	}
}
