package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/meetext/backend/internal/config"
	"github.com/meetext/backend/pkg/apperr"
)

type Claims struct {
	UserID      uuid.UUID `json:"user_id"`
	SessionID   uuid.UUID `json:"session_id"`
	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`
	Role        string    `json:"role,omitempty"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTService struct {
	cfg config.JWTConfig
}

func NewJWTService(cfg config.JWTConfig) *JWTService {
	return &JWTService{cfg: cfg}
}

func (s *JWTService) IssueTokenPair(userID, sessionID uuid.UUID) (*TokenPair, error) {
	access, err := s.sign(userID, sessionID, s.cfg.AccessSecret, s.cfg.AccessTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := s.sign(userID, sessionID, s.cfg.RefreshSecret, s.cfg.RefreshTTL)
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *JWTService) ValidateAccessToken(tokenStr string) (*Claims, error) {
	return s.parse(tokenStr, s.cfg.AccessSecret)
}

func (s *JWTService) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	return s.parse(tokenStr, s.cfg.RefreshSecret)
}

func (s *JWTService) AccessTTL() time.Duration {
	return s.cfg.AccessTTL
}

func (s *JWTService) RefreshTTL() time.Duration {
	return s.cfg.RefreshTTL
}

func (s *JWTService) sign(userID, sessionID uuid.UUID, secret string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("jwt: sign: %w", err)
	}
	return signed, nil
}

func (s *JWTService) parse(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apperr.ErrTokenExpired
		}
		return nil, apperr.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, apperr.ErrTokenInvalid
	}
	return claims, nil
}
