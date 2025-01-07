package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/LuciusMortified/video-conv-bot/internal/ent"
	"github.com/LuciusMortified/video-conv-bot/pkg/logger"
	"github.com/LuciusMortified/video-conv-bot/pkg/urlcheck"
)

type Service interface {
	Convert(ctx context.Context, params ent.ConvertParams) ent.ConvertStateChan
}

type Controller struct {
	service Service
	bot     *tgbotapi.BotAPI
	cfg     Config
}

type Config struct {
	Token         string
	UpdateTimeout time.Duration
}

func (c Config) Validate() error {
	if c.Token == "" {
		return errors.New("token is required")
	}
	if c.UpdateTimeout < 0 {
		return errors.New("invalid update timeout")
	}
	return nil
}

func New(service Service, cfg Config) (*Controller, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("cfg.Validate: %w", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("tgbotapi.NewBotAPI: %w", err)
	}

	_ = tgbotapi.SetLogger(&BotLogger{})

	return &Controller{
		service: service,
		bot:     bot,
		cfg:     cfg,
	}, nil
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
				With("update", upd).
				Debug("received update")

			if upd.Message == nil {
				continue
			}

			if err := c.handleMessage(ctx, upd.Message); err != nil {
				logger.
					With("update", upd).
					With("error", err).
					Error("failed to handle update")
			}
		}
	}
}

func (c *Controller) handleMessage(ctx context.Context, msg *tgbotapi.Message) error {
	url, err := c.extractURL(msg)
	if err != nil {
		return fmt.Errorf("c.extractURL: %w", err)
	}

	if url == "" {
		return nil
	}

	stateChan := c.service.Convert(ctx, ent.ConvertParams{URL: url})
	if err = c.processReply(ctx, msg, stateChan); err != nil {
		return fmt.Errorf("c.processReply: %w", err)
	}

	return nil
}

func (c *Controller) processReply(ctx context.Context, msg *tgbotapi.Message, stateChan ent.ConvertStateChan) error {
	replyID := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case state, ok := <-stateChan:
			if !ok {
				return nil
			}

			err := c.handleConvertState(msg, state, &replyID)
			if err != nil {
				return fmt.Errorf("c.handleConvertState: %w", err)
			}
		}
	}
}

func (c *Controller) handleConvertState(msg *tgbotapi.Message, state ent.ConvertState, replyID *int) error {
	defer func() { _ = state.Cleanup() }()

	buf := new(bytes.Buffer)
	if err := ent.ConvertTemplates[state.Status].Execute(buf, state); err != nil {
		return fmt.Errorf("ent.ConvertTemplates[].Execute: %w", err)
	}
	replyText := buf.String()

	if *replyID == 0 {
		reply := tgbotapi.NewMessage(msg.Chat.ID, replyText)
		reply.ReplyToMessageID = msg.MessageID

		replyMsg, err := c.bot.Send(reply)
		if err != nil {
			logger.With("error", err).Error("Failed to send reply")
			return nil
		}

		*replyID = replyMsg.MessageID
	} else {
		_, err := c.bot.Request(tgbotapi.NewEditMessageText(msg.Chat.ID, *replyID, replyText))
		if err != nil {
			logger.With("error", err).Error("Failed to edit reply")
			return nil
		}
	}

	if state.Status == ent.ConvertDone {
		return c.handleConvertDone(msg, state, replyID)
	}

	return nil
}

func (c *Controller) handleConvertDone(msg *tgbotapi.Message, state ent.ConvertState, replyID *int) error {
	newReply := tgbotapi.NewVideo(msg.Chat.ID, tgbotapi.FileReader{
		Name:   state.Result.Filename,
		Reader: state.Result.Data,
	})
	newReply.ReplyToMessageID = msg.MessageID

	_, err := c.bot.Request(newReply)
	if err != nil {
		logger.With("error", err).Error("Failed to send video")
		return nil
	}

	_, err = c.bot.Request(tgbotapi.NewDeleteMessage(msg.Chat.ID, *replyID))
	if err != nil {
		logger.With("error", err).Error("Failed to delete newReply")
		return nil
	}

	return nil
}

func (c *Controller) extractURL(msg *tgbotapi.Message) (string, error) {
	var (
		url string
		err error
	)
	if msg.Video != nil {
		url, err = c.bot.GetFileDirectURL(msg.Video.FileID)
		if err != nil {
			return "", fmt.Errorf("c.bot.GetFileDirectURL: %w", err)
		}

		return url, nil
	}

	if msg.Document != nil {
		url, err = c.bot.GetFileDirectURL(msg.Document.FileID)
		if err != nil {
			return "", fmt.Errorf("c.bot.GetFileDirectURL: %w", err)
		}

		return url, nil
	}

	if msg.Text != "" && urlcheck.IsUrl(msg.Text) {
		url = msg.Text

		return url, nil
	}

	return url, err
}
