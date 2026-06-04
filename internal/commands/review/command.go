package review

import (
	"fmt"

	"mentat/internal/context"
	"mentat/internal/core"
	tuireview "mentat/internal/tui/review"

	"github.com/gdamore/tcell/v2"
	"github.com/open-spaced-repetition/go-fsrs/v3"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "review",
		Short: "Review due cards",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := context.GetMentatContext(cmd.Context())
			if err != nil {
				return err
			}

			if core.CardsDue(m.Decks) == 0 {
				fmt.Println("no cards due")
				return nil
			}

			app := tview.NewApplication()
			view := tuireview.NewView()
			cr := tuireview.NewCardReview(m.History, m.Decks, view)
			if cr == nil {
				fmt.Println("no cards due")
				return nil
			}

			view.Text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Rune() == 'q' {
					app.Stop()
					return nil
				}
				if cr.HasReveal {
					if event.Rune() == ' ' {
						cr.Reveal()
						return nil
					}
				} else {
					if rating, ok := runeToRating(event.Rune()); ok {
						if !cr.Next(rating) {
							app.Stop()
						}
						return nil
					}
				}
				return event
			})

			app.SetRoot(view.Flex, true)
			return app.Run()
		},
	}
}

func runeToRating(r rune) (fsrs.Rating, bool) {
	switch r {
	case '1':
		return fsrs.Again, true
	case '2':
		return fsrs.Hard, true
	case '3':
		return fsrs.Good, true
	case '4':
		return fsrs.Easy, true
	default:
		return 0, false
	}
}
