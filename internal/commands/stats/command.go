package stats

import (
	"fmt"

	"mentat/internal/context"

	"github.com/open-spaced-repetition/go-fsrs/v3"
	"github.com/spf13/cobra"
)

func stateToString(state fsrs.State) string {
	switch state {
	case fsrs.New:
		return "New"
	case fsrs.Learning:
		return "Learning"
	case fsrs.Relearning:
		return "Relearning"
	case fsrs.Review:
		return "Review"
	}
	return ""
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Report Statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := context.GetMentatContext(cmd.Context())
			if err != nil {
				return err
			}

			fmt.Printf("%d decks found\n", len(m.Decks))
			for i, d := range m.Decks {
				stats := d.Stats()

				fmt.Printf(
					"%d) name: %s\n   cards: %d\n   due: %d\n   state-new: %d\n   state-learning: %d\n   state-relearning: %d\n   state-review: %d\n   last reviewed: %s\n",
					i+1,
					d.Name,
					stats.Cards.Count,
					stats.Cards.Due,
					stats.Cards.New,
					stats.Cards.Learning,
					stats.Cards.Relearning,
					stats.Cards.Review,
					stats.LastReview.Local().Format("2006-01-02"),
				)
			}

			return nil
		},
	}
	return cmd
}
