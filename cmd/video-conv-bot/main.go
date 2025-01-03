package main

import (
	"github.com/LuciusMortified/video-conv-bot/cmd/video-conv-bot/cmd"
	"github.com/LuciusMortified/video-conv-bot/pkg/logger"
)

func main() {
	defer logger.Flush()

	if err := cmd.Execute(); err != nil {
		logger.Error("Failed to execute", logger.NewField("error", err))
	}
}
