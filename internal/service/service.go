package service

import (
	"context"
	"github.com/LuciusMortified/video-conv-bot/internal/ent"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) Convert(ctx context.Context, params ent.ConvertParams) error {
	return nil
}
