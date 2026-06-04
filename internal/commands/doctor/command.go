package doctor

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	return &cobra.Command{

		Use: "doctor",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
