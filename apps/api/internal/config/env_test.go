package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadReadsOAuthConfigFromEnvFile(t *testing.T) {
	clearConfigEnv(t)

	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := strings.Join([]string{
		"DATABASE_URL=postgres://meetext:meetext@localhost:5432/meetext?sslmode=disable",
		"JWT_ACCESS_SECRET=12345678901234567890123456789012",
		"JWT_REFRESH_SECRET=abcdefghijklmnopqrstuvwxyz123456",
		"GOOGLE_CLIENT_ID=test-google-client-id",
		"GOOGLE_CLIENT_SECRET=test-google-client-secret",
		"GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/oauth/google/callback",
		"OAUTH_STATE_SECRET=state-secret-with-at-least-thirty-two-chars",
	}, "\n")

	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	cfg, err := Load(envPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.OAuth.GoogleClientID != "test-google-client-id" {
		t.Fatalf("GoogleClientID = %q", cfg.OAuth.GoogleClientID)
	}
	if cfg.OAuth.GoogleClientSecret != "test-google-client-secret" {
		t.Fatalf("GoogleClientSecret = %q", cfg.OAuth.GoogleClientSecret)
	}
	if cfg.OAuth.GoogleRedirectURL != "http://localhost:8080/api/v1/auth/oauth/google/callback" {
		t.Fatalf("GoogleRedirectURL = %q", cfg.OAuth.GoogleRedirectURL)
	}
	if err := cfg.ValidateAPI(); err != nil {
		t.Fatalf("validate api config: %v", err)
	}
}

func TestLoadUsesEnvironmentOverrideForGoogleClientID(t *testing.T) {
	clearConfigEnv(t)

	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := strings.Join([]string{
		"GOOGLE_CLIENT_ID=file-google-client-id",
		"GOOGLE_CLIENT_SECRET=test-google-client-secret",
		"GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/oauth/google/callback",
		"OAUTH_STATE_SECRET=state-secret-with-at-least-thirty-two-chars",
	}, "\n")

	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	t.Setenv("GOOGLE_CLIENT_ID", "env-google-client-id")

	cfg, err := Load(envPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.OAuth.GoogleClientID != "env-google-client-id" {
		t.Fatalf("GoogleClientID = %q", cfg.OAuth.GoogleClientID)
	}
}

func TestValidateAPIAllowsPasswordAuthWhenGoogleOAuthIsMissing(t *testing.T) {
	cfg := &Config{
		App: AppConfig{FrontendURL: "http://localhost:3000"},
		DB:  DBConfig{DSN: "postgres://meetext:meetext@localhost:5432/meetext?sslmode=disable"},
		JWT: JWTConfig{
			AccessSecret:  "12345678901234567890123456789012",
			RefreshSecret: "abcdefghijklmnopqrstuvwxyz123456",
		},
		OAuth: OAuthConfig{
			GoogleRedirectURL: "http://localhost:8080/api/v1/auth/oauth/google/callback",
			StateSecret:       "state-secret-with-at-least-thirty-two-chars",
		},
	}

	if err := cfg.ValidateAPI(); err != nil {
		t.Fatalf("validate api config: %v", err)
	}
}

func clearConfigEnv(t *testing.T) {
	t.Helper()

	for _, key := range []string{
		"DATABASE_URL",
		"JWT_ACCESS_SECRET",
		"JWT_REFRESH_SECRET",
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		"GOOGLE_REDIRECT_URL",
		"OAUTH_STATE_SECRET",
	} {
		t.Setenv(key, "")
	}
}
