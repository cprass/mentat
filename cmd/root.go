package cmd

import (
	"fmt"
	"os"

	"mentat/internal/commands/doctor"
	"mentat/internal/commands/review"
	"mentat/internal/commands/run"
	"mentat/internal/commands/stats"
	"mentat/internal/config"
	"mentat/internal/context"
	"mentat/internal/core"
	"mentat/internal/sync"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:                "mentat",
		PersistentPreRunE:  onInit,
		PersistentPostRunE: onExit,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var configFile string

func init() {
	// This sets up flags and defaults as well as viper/cobra bindings for all app config options
	config.InitFlags(rootCmd)

	rootCmd.AddCommand(stats.NewCommand())
	rootCmd.AddCommand(review.NewCommand())
	rootCmd.AddCommand(run.NewCommand())
	rootCmd.AddCommand(doctor.NewCommand())
}

func onInit(cmd *cobra.Command, args []string) error {
	m, err := core.NewMentat()
	if err != nil {
		return err
	}
	ctx := context.SetMentatContext(cmd.Context(), m)
	cmd.SetContext(ctx)

	if config.SyncEnabled() {
		if err := sync.SyncPull(ctx); err != nil {
			return err
		}
	}
	return nil
}

func onExit(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	mentat, err := context.GetMentatContext(ctx)
	if err != nil {
		return err
	}
	mentat.Close()

	if config.SyncEnabled() {
		if err := sync.SyncPush(ctx); err != nil {
			return err
		}
	}
	return nil
}
