package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/LuciusMortified/video-conv-bot/internal/config"
	"github.com/LuciusMortified/video-conv-bot/pkg/logger"
)

var (
	cfg config.Config

	rootCmd = &cobra.Command{
		Use:   "video-conv-bot",
		Short: "Telegram bot than converting video files to MP4 format.",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(runCmd)
}

func initConfig() {
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		logger.
			With("error", err).
			Fatal("Failed to unmarshal config")
	}
}
