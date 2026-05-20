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
		Options: map[string]interface{}{
			"num_ctx": 131072, // Force massive 128k context window for long transcripts
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", p.baseURL)
	
	maxRetries := 3
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		p.log.Debug().Str("url", url).Str("model", p.model).Int("attempt", attempt).Msg("sending request to ollama")
		
		// Enforce a strict 3600s timeout specifically for massive LLM inference calls
		reqCtx, cancel := context.WithTimeout(ctx, 3600*time.Second)
		req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			cancel()
			return "", fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := p.client.Do(req)
		if err != nil {
			cancel()
			lastErr = fmt.Errorf("failed to execute request on attempt %d: %w", attempt, err)
			p.log.Warn().Err(lastErr).Msg("ollama request failed, retrying...")
			time.Sleep(2 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			cancel()
			lastErr = fmt.Errorf("ollama API error: status %d, response: %s", resp.StatusCode, string(bodyBytes))
			p.log.Warn().Err(lastErr).Msg("ollama returned non-200 status, retrying...")
			time.Sleep(2 * time.Second)
			continue
		}

		var genResp generateResponse
		err = json.NewDecoder(resp.Body).Decode(&genResp)
		resp.Body.Close()
		cancel()

		if err != nil {
			lastErr = fmt.Errorf("failed to decode response on attempt %d: %w", attempt, err)
			p.log.Warn().Err(lastErr).Msg("failed to decode ollama response, retrying...")
			time.Sleep(2 * time.Second)
			continue
		}

		cleanJSON := sanitizeJSON(genResp.Response)
		
		// Structured parsing validation: verify the LLM gave us valid JSON.
		if !json.Valid([]byte(cleanJSON)) {
			lastErr = fmt.Errorf("ollama returned invalid json on attempt %d", attempt)
			p.log.Warn().Err(lastErr).Msg("invalid json structure, retrying...")
			time.Sleep(2 * time.Second)
			continue
		}

		p.log.Info().Msg("ollama generation successful")
		return cleanJSON, nil
	}

	return "", fmt.Errorf("failed to generate JSON after %d attempts. Last error: %w", maxRetries, lastErr)
}

// sanitizeJSON acts as a hallucination protection layer, trimming markdown block backticks that LLMs often incorrectly inject.
func sanitizeJSON(input string) string {
	s := strings.TrimSpace(input)
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
	}
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSuffix(s, "```")
	}
	return strings.TrimSpace(s)
}
