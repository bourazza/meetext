package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	infraoauth "github.com/meetext/backend/internal/infrastructure/oauth"
	ucauth "github.com/meetext/backend/internal/usecase/auth"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/meetext/backend/pkg/response"
	"github.com/rs/zerolog"
)

const oauthStateCookie = "oauth_state"

// OAuthHandler handles the redirect and callback for a single OAuth provider.
type OAuthHandler struct {
	provider    *infraoauth.Provider
	authUC      *ucauth.UseCase
	frontendURL string
	log         zerolog.Logger
}

func NewOAuthHandler(provider *infraoauth.Provider, authUC *ucauth.UseCase, frontendURL string, log zerolog.Logger) *OAuthHandler {
	return &OAuthHandler{
		provider:    provider,
		authUC:      authUC,
		frontendURL: frontendURL,
		log:         log.With().Str("provider", string(provider.Name())).Logger(),
	}
}

// Redirect generates a signed state, stores it in a short-lived cookie, and
// redirects the browser to the provider's consent screen.
//
// GET /api/v1/auth/oauth/{provider}
func (h *OAuthHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	if !h.provider.IsConfigured() {
		h.log.Error().Msg("oauth: provider not configured — client_id or client_secret is empty")
		response.Error(w, apperr.Wrap(nil, http.StatusServiceUnavailable, "OAUTH_NOT_CONFIGURED", "OAuth provider is not configured"))
		return
	}

	nonce, err := randomNonce()
	if err != nil {
		h.log.Error().Err(err).Msg("oauth: failed to generate nonce")
		response.Error(w, apperr.ErrInternal)
		return
	}

	state := h.provider.GenerateState(nonce)

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})

	authURL := h.provider.AuthURL(state)
	h.log.Debug().Str("redirect_to", authURL).Msg("oauth: redirecting to provider")
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Callback handles the provider redirect, exchanges the code, creates a server
// session, sets secure auth cookies, and redirects the browser to the dashboard.
//
// GET /api/v1/auth/oauth/{provider}/callback
func (h *OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	// Check for provider-side error first
	if errParam := r.URL.Query().Get("error"); errParam != "" {
		h.log.Warn().Str("provider_error", errParam).
			Str("description", r.URL.Query().Get("error_description")).
			Msg("oauth: provider returned error")
		h.redirectError(w, r, errParam)
		return
	}

	// Validate state cookie
	cookie, err := r.Cookie(oauthStateCookie)
	if err != nil || cookie.Value == "" {
		h.log.Warn().Err(err).Msg("oauth: missing state cookie")
		h.redirectError(w, r, "missing_state")
		return
	}

	// Clear the cookie immediately
	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	stateParam := r.URL.Query().Get("state")

	if !h.provider.ValidateState(stateParam) {
		h.log.Warn().Str("state", stateParam).Msg("oauth: invalid state HMAC")
		h.redirectError(w, r, "invalid_state")
		return
	}
	if cookie.Value != stateParam {
		h.log.Warn().Msg("oauth: state cookie mismatch")
		h.redirectError(w, r, "state_mismatch")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		h.log.Warn().Msg("oauth: missing code parameter")
		h.redirectError(w, r, "missing_code")
		return
	}

	// Exchange code → user info
	info, err := h.provider.Exchange(r.Context(), code)
	if err != nil {
		h.log.Error().Err(err).Msg("oauth: code exchange failed")
		h.redirectError(w, r, "exchange_failed")
		return
	}

	h.log.Debug().Str("email", info.Email).Str("provider_id", info.ProviderID).Msg("oauth: exchange successful")

	// Find or create user, issue JWT
	res, err := h.authUC.OAuthLogin(r.Context(), h.provider.Name(), info, requestMeta(r))
	if err != nil {
		h.log.Error().Err(err).Str("email", info.Email).Msg("oauth: login failed")
		h.redirectError(w, r, "login_failed")
		return
	}

	h.log.Info().Str("user_id", res.User.ID.String()).Str("email", res.User.Email).Msg("oauth: login successful")
	setAuthCookies(w, r, res.AccessToken, res.RefreshToken, true)

	redirectURL := fmt.Sprintf("%s/auth/callback?success=true", h.frontendURL)
	h.log.Debug().Str("redirect_to", redirectURL).Msg("oauth: redirecting to frontend callback")
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) redirectError(w http.ResponseWriter, r *http.Request, reason string) {
	redirectURL := fmt.Sprintf("%s/login?error=%s", h.frontendURL, url.QueryEscape(reason))
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func randomNonce() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
