package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	authdomain "github.com/meetext/backend/internal/domain/auth"
	infraauth "github.com/meetext/backend/internal/infrastructure/auth"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/meetext/backend/pkg/constants"
	"github.com/meetext/backend/pkg/response"
	"github.com/rs/zerolog"
)

func Auth(jwt *infraauth.JWTService, sessions authdomain.TokenRepository, log zerolog.Logger) func(http.Handler) http.Handler {
	authLog := log.With().Str("component", "auth_middleware").Logger()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := bearerToken(r)
			if token == "" {
				authLog.Debug().Str("path", r.URL.Path).Msg("auth: missing access token")
				response.Error(w, apperr.ErrUnauthorized)
				return
			}

			claims, err := jwt.ValidateAccessToken(token)
			if err != nil {
				authLog.Debug().Err(err).Str("path", r.URL.Path).Msg("auth: invalid access token")
				response.Error(w, err)
				return
			}
			session, err := sessions.GetSession(r.Context(), claims.SessionID)
			if err != nil || session.UserID != claims.UserID || session.RevokedAt != nil || time.Now().After(session.ExpiresAt) {
				authLog.Debug().Err(err).Str("user_id", claims.UserID.String()).Str("session_id", claims.SessionID.String()).Msg("auth: session validation failed")
				response.Error(w, apperr.ErrUnauthorized)
				return
			}

			authLog.Debug().Str("user_id", claims.UserID.String()).Str("session_id", claims.SessionID.String()).Msg("auth: session validated")
			ctx := context.WithValue(r.Context(), constants.CtxUserID, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	cookie, err := r.Cookie("meetext_access")
	if err != nil {
		return ""
	}
	return cookie.Value
}
