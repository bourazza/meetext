package middleware

import (
	"context"
	"net/http"
	"strings"

	infraauth "github.com/meetext/backend/internal/infrastructure/auth"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/meetext/backend/pkg/constants"
	"github.com/meetext/backend/pkg/response"
)

func Auth(jwt *infraauth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				response.Error(w, apperr.ErrUnauthorized)
				return
			}

			token := strings.TrimPrefix(header, "Bearer ")
			claims, err := jwt.ValidateAccessToken(token)
			if err != nil {
				response.Error(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), constants.CtxUserID, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
