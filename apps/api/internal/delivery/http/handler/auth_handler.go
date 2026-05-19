package handler

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	ucauth "github.com/meetext/backend/internal/usecase/auth"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/meetext/backend/pkg/constants"
	"github.com/meetext/backend/pkg/response"
	"github.com/meetext/backend/pkg/validator"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	uc  *ucauth.UseCase
	log zerolog.Logger
}

func NewAuthHandler(uc *ucauth.UseCase, log zerolog.Logger) *AuthHandler {
	return &AuthHandler{uc: uc, log: log.With().Str("component", "auth_handler").Logger()}
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in ucauth.RegisterInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	res, err := h.uc.Register(r.Context(), in, requestMeta(r))
	if err != nil {
		h.log.Warn().Err(err).Str("email", in.Email).Msg("auth: registration failed")
		response.Error(w, err)
		return
	}
	setAuthCookies(w, r, res.AccessToken, res.RefreshToken, true)
	h.log.Info().Str("user_id", res.User.ID.String()).Str("email", res.User.Email).Msg("auth: user registered and session cookie issued")
	response.Created(w, res)
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in ucauth.LoginInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	meta := requestMeta(r)
	meta.Remember = in.RememberMe
	res, err := h.uc.Login(r.Context(), in, meta)
	if err != nil {
		h.log.Warn().Err(err).Str("email", in.Email).Msg("auth: login failed")
		response.Error(w, err)
		return
	}
	setAuthCookies(w, r, res.AccessToken, res.RefreshToken, in.RememberMe)
	h.log.Info().Str("user_id", res.User.ID.String()).Str("email", res.User.Email).Msg("auth: login succeeded and session cookie issued")
	response.OK(w, res)
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if r.Body != nil {
		_, _ = validator.Decode(r, &body)
	}
	if body.RefreshToken == "" {
		body.RefreshToken = readCookie(r, refreshCookieName)
	}

	tokens, err := h.uc.RefreshToken(r.Context(), body.RefreshToken)
	if err != nil {
		h.log.Warn().Err(err).Msg("auth: refresh failed")
		response.Error(w, err)
		return
	}
	setAuthCookies(w, r, tokens.AccessToken, tokens.RefreshToken, true)
	h.log.Debug().Msg("auth: refresh succeeded and cookies rotated")
	response.OK(w, tokens)
}

// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_ = h.uc.Logout(r.Context(), readCookie(r, refreshCookieName))
	clearAuthCookies(w, r)
	h.log.Info().Msg("auth: logout completed and cookies cleared")
	response.OK(w, map[string]string{"message": "Signed out."})
}

// POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var in ucauth.ForgotPasswordInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.uc.RequestPasswordReset(r.Context(), in); err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, map[string]string{"message": "If that email exists, a reset link has been sent."})
}

// POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var in ucauth.ResetPasswordInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.uc.ResetPassword(r.Context(), in); err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, map[string]string{"message": "Your password has been updated."})
}

// POST /api/v1/auth/verify-email
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var in ucauth.VerifyEmailInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.uc.VerifyEmail(r.Context(), in); err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, map[string]string{"message": "Email verified."})
}

// POST /api/v1/auth/resend-verification
func (h *AuthHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var in ucauth.ResendVerificationInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.uc.ResendVerification(r.Context(), in); err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, map[string]string{"message": "If that account needs verification, a new link has been sent."})
}

// GET /api/v1/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constants.CtxUserID).(uuid.UUID)
	if !ok {
		response.Error(w, apperr.ErrUnauthorized)
		return
	}

	u, err := h.uc.CurrentUser(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, u)
}

const (
	accessCookieName  = "meetext_access"
	refreshCookieName = "meetext_refresh"
)

func setAuthCookies(w http.ResponseWriter, r *http.Request, access, refresh string, remember bool) {
	secure := isSecureRequest(r)
	accessMaxAge := int((15 * time.Minute).Seconds())
	refreshMaxAge := int((7 * 24 * time.Hour).Seconds())
	if !remember {
		refreshMaxAge = int((24 * time.Hour).Seconds())
	}

	http.SetCookie(w, &http.Cookie{
		Name:     accessCookieName,
		Value:    access,
		Path:     "/",
		MaxAge:   accessMaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    refresh,
		Path:     "/",
		MaxAge:   refreshMaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	})
}

func clearAuthCookies(w http.ResponseWriter, r *http.Request) {
	secure := isSecureRequest(r)
	for _, name := range []string{accessCookieName, refreshCookieName} {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   secure,
		})
	}
}

func requestMeta(r *http.Request) ucauth.RequestMeta {
	return ucauth.RequestMeta{
		UserAgent: r.UserAgent(),
		IP:        clientIP(r),
		Remember:  true,
	}
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	host := r.RemoteAddr
	if ip, _, err := net.SplitHostPort(host); err == nil {
		return strings.Trim(ip, "[]")
	}
	return strings.Trim(host, "[]")
}

func readCookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func isSecureRequest(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}
