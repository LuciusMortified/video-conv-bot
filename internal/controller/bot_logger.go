package controller

import (
	"fmt"
	"github.com/LuciusMortified/video-conv-bot/pkg/logger"
)

type BotLogger struct{}

func (l BotLogger) Println(v ...interface{}) {
	logger.Info(fmt.Sprintf("telegram bot: %s", fmt.Sprint(v...)))
}

func (l BotLogger) Printf(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf("telegram bot: %s", fmt.Sprintf(format, v...)))
}
