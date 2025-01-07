package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/LuciusMortified/video-conv-bot/internal/controller"
	"github.com/LuciusMortified/video-conv-bot/internal/service"
	"github.com/LuciusMortified/video-conv-bot/pkg/logger"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run telegram bot",
		RunE:  run,
	}
)

func init() {
	pFlags := runCmd.PersistentFlags()
	pFlags.StringVar(&cfg.Telegram.Token, "token", "", "Telegram bot token")
	_ = viper.BindPFlag("telegram.token", pFlags.Lookup("token"))
	_ = viper.BindEnv("telegram.token", "TELEGRAM_TOKEN")

	pFlags.StringVar(&cfg.Convert.StoragePath, "storage-path", "", "Storage path for converting files")
	_ = viper.BindPFlag("convert.storage_path", pFlags.Lookup("storage-path"))
	_ = viper.BindEnv("convert.storage_path", "CONVERT_STORAGE_PATH")
}

func run(cmd *cobra.Command, _ []string) error {
	logger.Info("Starting telegram bot")

	svc, err := service.New(service.Config{
		StoragePath: cfg.Convert.StoragePath,
	})
	if err != nil {
		return fmt.Errorf("service.New: %w", err)
	}

	ctrl, err := controller.New(svc, controller.Config{
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
