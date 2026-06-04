package router

import (
	"mentat/internal/commands/run/stats"
	"mentat/internal/core"
	tuireview "mentat/internal/tui/review"

	"github.com/gdamore/tcell/v2"
	"github.com/open-spaced-repetition/go-fsrs/v3"
	"github.com/rivo/tview"
)

type router struct {
	mentat *core.Mentat
	app    *tview.Application
}

func (r *router) ShowStats() {
	stats.RenderStats(r.mentat, r)
}

func (r *router) ShowReview(decks []*core.Deck) {
	view := tuireview.NewView()
	cr := tuireview.NewCardReview(r.mentat.History, decks, view)
	if cr == nil {
		r.ShowStats()
		return
	}

	view.Text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			r.ShowStats()
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
					r.ShowStats()
				}
				return nil
			}
		}
		return event
	})

	r.app.SetRoot(view.Flex, true)
}

func (r *router) RenderView(view tview.Primitive) {
	r.app.SetRoot(view, true)
}

func (r *router) App() *tview.Application {
	return r.app
}

func NewRouter(app *tview.Application, mentat *core.Mentat) *router {
	return &router{app: app, mentat: mentat}
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
