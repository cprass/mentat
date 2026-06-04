package history

import (
	"fmt"
	"time"

	"github.com/open-spaced-repetition/go-fsrs/v3"
)

type ReviewEvent struct {
	CardID string
	Rating fsrs.Rating
	Card   fsrs.Card
}

type Store interface {
	Append(event *ReviewEvent) error
	LoadAll() ([]*ReviewEvent, error)
	Close() error
}

func (re *ReviewEvent) String() string {
	return fmt.Sprintf("%s %d %s %f %f %d %d %d %d %d %s",
		re.CardID,
		re.Rating,
		re.Card.Due.Format(time.RFC3339),
		re.Card.Stability,
		re.Card.Difficulty,
		re.Card.ElapsedDays,
		re.Card.ScheduledDays,
		re.Card.Reps,
		re.Card.Lapses,
		re.Card.State,
		re.Card.LastReview.Format(time.RFC3339),
	)
}

func NewReviewEvent(cardID string, rating fsrs.Rating, card fsrs.Card) *ReviewEvent {
	return &ReviewEvent{
		CardID: cardID,
		Rating: rating,
		Card:   card,
	}
}

func NewReviewEventFromString(raw string) (*ReviewEvent, error) {
	var (
		cardID, due, lastReview          string
		rating, state                    int
		stability, difficulty            float64
		elapsed, scheduled, reps, lapses uint64
	)

	_, err := fmt.Sscanf(raw, "%s %d %s %f %f %d %d %d %d %d %s",
		&cardID, &rating, &due, &stability, &difficulty, &elapsed, &scheduled, &reps, &lapses, &state, &lastReview,
	)
	if err != nil {
		return nil, err
	}

	dueTime, err := time.Parse(time.RFC3339, due)
	if err != nil {
		return nil, err
	}

	lastReviewTime, err := time.Parse(time.RFC3339, lastReview)
	if err != nil {
		return nil, err
	}

	return &ReviewEvent{
		CardID: cardID,
		Rating: fsrs.Rating(rating),
		Card: fsrs.Card{
			Due:           dueTime,
			Stability:     stability,
			Difficulty:    difficulty,
			ElapsedDays:   elapsed,
			ScheduledDays: scheduled,
			Reps:          reps,
			Lapses:        lapses,
			State:         fsrs.State(state),
			LastReview:    lastReviewTime,
		},
	}, nil
}
