package cmd

import (
	"fmt"
	"github.com/LuciusMortified/video-conv-bot/internal/controller"
	"github.com/LuciusMortified/video-conv-bot/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run telegram bot",
		RunE:  run,
	}
)

func run(cmd *cobra.Command, _ []string) error {
	logger.Info("Starting telegram bot")

	ctrl, err := controller.New(controller.Config{
		Token: cfg.Telegram.Token,
	})
	if err != nil {
		return fmt.Errorf("controller.New: %w", err)
	}

	if err = ctrl.Run(cmd.Context()); err != nil {
		return fmt.Errorf("ctrl.Run: %w", err)
	}

	return nil
}
