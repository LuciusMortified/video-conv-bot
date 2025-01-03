package cmd

import (
	"github.com/LuciusMortified/video-conv-bot/internal/config"
	"github.com/LuciusMortified/video-conv-bot/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     config.Config

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

	pFlags := rootCmd.PersistentFlags()
	pFlags.StringVar(&cfgFile, "config", "", "config file path")
	_ = viper.BindPFlag("config", pFlags.Lookup("config"))

	rootCmd.AddCommand(runCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.
			WithField("file", viper.ConfigFileUsed()).
			WithField("error", err).
			Fatal("Failed to read config file")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logger.
			WithField("file", viper.ConfigFileUsed()).
			WithField("error", err).
			Fatal("Failed to unmarshal config file")
	}
}
