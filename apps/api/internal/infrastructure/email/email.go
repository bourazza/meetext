package email

import (
	"context"

	"github.com/rs/zerolog"
)

type Service interface {
	SendVerification(ctx context.Context, to, name, link string) error
	SendPasswordReset(ctx context.Context, to, name, link string) error
}

type LogService struct {
	log zerolog.Logger
}

func NewLogService(log zerolog.Logger) *LogService {
	return &LogService{log: log.With().Str("component", "email").Logger()}
}

func (s *LogService) SendVerification(ctx context.Context, to, name, link string) error {
	s.log.Info().Str("to", to).Str("name", name).Str("link", link).Msg("verification email queued")
	return nil
}

func (s *LogService) SendPasswordReset(ctx context.Context, to, name, link string) error {
	s.log.Info().Str("to", to).Str("name", name).Str("link", link).Msg("password reset email queued")
	return nil
}
