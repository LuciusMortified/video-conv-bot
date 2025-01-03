package controller

import (
	"context"
	"fmt"
	"github.com/LuciusMortified/video-conv-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type Controller struct {
	bot *tgbotapi.BotAPI
	cfg Config
}

type Config struct {
	Token         string
	Debug         bool
	UpdateTimeout time.Duration
}

func New(cfg Config) (*Controller, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("tgbotapi.NewBotAPI: %w", err)
	}

	bot.Debug = cfg.Debug
	_ = tgbotapi.SetLogger(&BotLogger{})

	return &Controller{bot: bot, cfg: cfg}, nil
}

func (c *Controller) Run(ctx context.Context) error {
	updConfig := tgbotapi.NewUpdate(0)
	updConfig.Timeout = int(c.cfg.UpdateTimeout.Seconds())

	updChan := c.bot.GetUpdatesChan(updConfig)
	logger.Info("Start receiving updates")

	for {
		select {
		case <-ctx.Done():
			c.bot.StopReceivingUpdates()
			logger.Info("Stop receiving updates")

		case upd, ok := <-updChan:
			if !ok {
				return nil
			}

			logger.
				WithField("update", upd).
				Debug("received update")

			if err := c.HandleUpdate(ctx, upd); err != nil {
				logger.
					WithField("update", upd).
					WithField("error", err).
					Error("failed to handle update")
			}
		}
	}
}

func (c *Controller) HandleUpdate(ctx context.Context, upd tgbotapi.Update) error {

	return nil
}
