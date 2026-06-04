package core

import "time"

type CardStats struct {
	Count      int
	Due        int
	New        int
	Learning   int
	Relearning int
	Review     int
}

type DeckStats struct {
	Cards      CardStats
	LastReview time.Time
}

func (s DeckStats) Add(s2 DeckStats) DeckStats {
	latestReview := s.LastReview
	if latestReview.Before(s2.LastReview) {
		latestReview = s2.LastReview
	}

	return DeckStats{
		Cards: CardStats{
			Count:      s.Cards.Count + s2.Cards.Count,
			Due:        s.Cards.Due + s2.Cards.Due,
			New:        s.Cards.New + s2.Cards.New,
			Learning:   s.Cards.Learning + s2.Cards.Learning,
			Relearning: s.Cards.Relearning + s2.Cards.Relearning,
			Review:     s.Cards.Review + s2.Cards.Review,
		},
		LastReview: latestReview,
	}
}
