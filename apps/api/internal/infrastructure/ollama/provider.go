package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/meetext/backend/internal/config"
	"github.com/rs/zerolog"
)

const (
	defaultTimeout = 5 * time.Minute
	maxRetries     = 2
)

type Provider struct {
	baseURL string
	model   string
	client  *http.Client
	log     zerolog.Logger
}

type generateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Format  string                 `json:"format"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type generateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func NewProvider(cfg config.AIConfig, log zerolog.Logger) *Provider {
	baseURL := strings.TrimRight(cfg.OllamaURL, "/")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	model := cfg.OllamaModel
	if model == "" {
		model = "llama3.2:1b"
	}
	return &Provider{
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{Timeout: defaultTimeout + 10*time.Second},
		log:     log.With().Str("component", "ollama_provider").Logger(),
	}
}

// GenerateJSON sends a prompt to Ollama and returns a validated JSON string.
func (p *Provider) GenerateJSON(ctx context.Context, prompt string) (string, error) {
	reqBody := generateRequest{
		Model:  p.model,
		Prompt: prompt,
		Format: "json",
		Stream: false,
		Options: map[string]interface{}{
			"num_ctx":     4096,
			"temperature": 0.1,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("ollama: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", p.baseURL)
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		reqCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
		req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			cancel()
			return "", fmt.Errorf("ollama: create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		p.log.Info().Str("model", p.model).Int("attempt", attempt).Msg("ollama: sending request")

		resp, err := p.client.Do(req)
		cancel()
		if err != nil {
			lastErr = fmt.Errorf("ollama: request attempt %d: %w", attempt, err)
			p.log.Warn().Err(lastErr).Msg("ollama: request failed, retrying")
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("ollama: read body attempt %d: %w", attempt, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("ollama: status %d: %s", resp.StatusCode, string(body))
			p.log.Warn().Err(lastErr).Msg("ollama: non-200 response")
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
			continue
		}

		var genResp generateResponse
		if err := json.Unmarshal(body, &genResp); err != nil {
			lastErr = fmt.Errorf("ollama: decode response attempt %d: %w", attempt, err)
			continue
		}

		clean := sanitizeJSON(genResp.Response)
		if !json.Valid([]byte(clean)) {
			lastErr = fmt.Errorf("ollama: invalid JSON in response attempt %d", attempt)
			p.log.Warn().Str("raw", clean[:min(200, len(clean))]).Msg("ollama: invalid json, retrying")
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
			continue
		}

		p.log.Info().Msg("ollama: generation successful")
		return clean, nil
	}

	return "", fmt.Errorf("ollama: failed after %d attempts: %w", maxRetries, lastErr)
}

func sanitizeJSON(input string) string {
	s := strings.TrimSpace(input)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
