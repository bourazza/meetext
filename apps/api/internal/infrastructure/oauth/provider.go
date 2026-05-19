package oauth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/meetext/backend/internal/config"
	"github.com/meetext/backend/internal/domain/user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// UserInfo is the normalised profile returned by any provider.
type UserInfo struct {
	ProviderID string
	Email      string
	Name       string
	AvatarURL  string
}

// Provider wraps an oauth2.Config and knows how to fetch a UserInfo.
type Provider struct {
	name        user.Provider
	cfg         *oauth2.Config
	fetchUser   func(ctx context.Context, token *oauth2.Token) (*UserInfo, error)
	stateSecret string
}

// NewGoogle returns a configured Google OAuth provider.
func NewGoogle(cfg config.OAuthConfig) *Provider {
	return &Provider{
		name: user.ProviderGoogle,
		cfg: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
		fetchUser:   fetchGoogleUser,
		stateSecret: cfg.StateSecret,
	}
}

// NewGitHub returns a configured GitHub OAuth provider.
func NewGitHub(cfg config.OAuthConfig) *Provider {
	return &Provider{
		name: user.ProviderGitHub,
		cfg: &oauth2.Config{
			ClientID:     cfg.GitHubClientID,
			ClientSecret: cfg.GitHubClientSecret,
			RedirectURL:  cfg.GitHubRedirectURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
		fetchUser:   fetchGitHubUser,
		stateSecret: cfg.StateSecret,
	}
}

// Name returns the user.Provider constant for this provider.
func (p *Provider) Name() user.Provider { return p.name }

// IsConfigured returns false when this provider cannot perform a complete
// OAuth authorization code flow.
func (p *Provider) IsConfigured() bool {
	return p != nil &&
		p.cfg != nil &&
		strings.TrimSpace(p.cfg.ClientID) != "" &&
		strings.TrimSpace(p.cfg.ClientSecret) != "" &&
		strings.TrimSpace(p.cfg.RedirectURL) != "" &&
		strings.TrimSpace(p.stateSecret) != ""
}

func (p *Provider) Validate() error {
	if p == nil || p.cfg == nil {
		return fmt.Errorf("provider config is nil")
	}

	var missing []string
	if strings.TrimSpace(p.cfg.ClientID) == "" {
		missing = append(missing, "client_id")
	}
	if strings.TrimSpace(p.cfg.ClientSecret) == "" {
		missing = append(missing, "client_secret")
	}
	if strings.TrimSpace(p.cfg.RedirectURL) == "" {
		missing = append(missing, "redirect_url")
	}
	if strings.TrimSpace(p.stateSecret) == "" {
		missing = append(missing, "state_secret")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing %s", strings.Join(missing, ", "))
	}

	return nil
}

// AuthURL returns the provider consent-screen URL with a signed state token.
func (p *Provider) AuthURL(state string) string {
	return p.cfg.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// Exchange trades the code for a token and returns the user's profile.
func (p *Provider) Exchange(ctx context.Context, code string) (*UserInfo, error) {
	token, err := p.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth: exchange: %w", err)
	}
	return p.fetchUser(ctx, token)
}

// GenerateState creates an HMAC-signed state string for CSRF protection.
func (p *Provider) GenerateState(nonce string) string {
	mac := hmac.New(sha256.New, []byte(p.stateSecret))
	mac.Write([]byte(nonce))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	return nonce + "." + sig
}

// ValidateState verifies the HMAC signature on the returned state.
func (p *Provider) ValidateState(state string) bool {
	idx := strings.LastIndex(state, ".")
	if idx < 0 {
		return false
	}
	nonce := state[:idx]
	expected := p.GenerateState(nonce)
	return hmac.Equal([]byte(state), []byte(expected))
}

// ── Google ────────────────────────────────────────────────────────────────────

type googleUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func fetchGoogleUser(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("google: userinfo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google: userinfo status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("google: read body: %w", err)
	}
	var g googleUserInfo
	if err := json.Unmarshal(body, &g); err != nil {
		return nil, fmt.Errorf("google: decode userinfo: %w", err)
	}
	if g.Email == "" {
		return nil, fmt.Errorf("google: no email in userinfo")
	}
	return &UserInfo{ProviderID: g.Sub, Email: g.Email, Name: g.Name, AvatarURL: g.Picture}, nil
}

// ── GitHub ────────────────────────────────────────────────────────────────────

type githubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type githubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func fetchGitHubUser(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("github: user request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github: user status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("github: read body: %w", err)
	}
	var g githubUserInfo
	if err := json.Unmarshal(body, &g); err != nil {
		return nil, fmt.Errorf("github: decode user: %w", err)
	}

	email := g.Email
	if email == "" {
		email, err = fetchGitHubPrimaryEmail(ctx, client)
		if err != nil {
			return nil, err
		}
	}
	if email == "" {
		return nil, fmt.Errorf("github: no verified primary email found")
	}

	name := g.Name
	if name == "" {
		name = g.Login
	}
	return &UserInfo{
		ProviderID: fmt.Sprintf("%d", g.ID),
		Email:      email,
		Name:       name,
		AvatarURL:  g.AvatarURL,
	}, nil
}

func fetchGitHubPrimaryEmail(ctx context.Context, client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", fmt.Errorf("github: emails request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("github: read emails: %w", err)
	}
	var emails []githubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", fmt.Errorf("github: decode emails: %w", err)
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", nil
}
