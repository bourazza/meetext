package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/meetext/backend/internal/config"
	"github.com/rs/zerolog"
)

type Provider struct {
	baseURL string
	model   string
	client  *http.Client
	log     zerolog.Logger
}

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Format string `json:"format"`
	Stream bool   `json:"stream"`
}

type generateResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func NewProvider(cfg config.AIConfig, log zerolog.Logger) *Provider {
	baseURL := strings.TrimRight(cfg.OllamaURL, "/")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	model := cfg.OllamaModel
	if model == "" {
		model = "llama3"
	}

	return &Provider{
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{},
		log:     log.With().Str("component", "ollama_provider").Logger(),
	}
}

// GenerateJSON sends a prompt to Ollama enforcing JSON output.
func (p *Provider) GenerateJSON(ctx context.Context, prompt string) (string, error) {
	reqBody := generateRequest{
		Model:  p.model,
		Prompt: prompt,
		Format: "json", // Enforce JSON output at the Ollama level
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", p.baseURL)
	p.log.Debug().Str("url", url).Str("model", p.model).Msg("sending request to ollama")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		p.log.Error().Err(err).Msg("ollama request failed")
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p.log.Error().Int("status_code", resp.StatusCode).Msg("ollama returned non-200 status")
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var genResp generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		p.log.Error().Err(err).Msg("failed to decode ollama response")
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	p.log.Info().Msg("ollama generation successful")
	return genResp.Response, nil
}
