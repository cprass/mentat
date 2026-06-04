package run

import (
	"mentat/internal/commands/run/router"
	"mentat/internal/context"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Start interactive review session",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := context.GetMentatContext(cmd.Context())
			if err != nil {
				return err
			}

			app := tview.NewApplication()

			r := router.NewRouter(app, m)

			r.ShowStats()

			return app.Run()
		},
	}
}
