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
	defaultTimeout = 10 * time.Minute // Longer timeout for 8B models
	maxRetries     = 2
	keepAlive      = "30m" // Keep model warm in memory
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
	Format  string                 `json:"format,omitempty"`
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
		model = "llama3.1:8b-instruct-q4_K_M"
	}
	return &Provider{
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{Timeout: defaultTimeout + 30*time.Second},
		log:     log.With().Str("component", "ollama_provider").Logger(),
	}
}

// GenerateJSON sends a prompt to Ollama with optimized settings for long-context processing.
func (p *Provider) GenerateJSON(ctx context.Context, prompt string) (string, error) {
	// Don't use "json" format - it causes null responses in some Ollama versions
	// Instead, rely on prompt instructions to return JSON
	return p.generate(ctx, prompt, "")
}

// GenerateText sends a prompt to Ollama for plain text generation.
func (p *Provider) GenerateText(ctx context.Context, prompt string) (string, error) {
	return p.generate(ctx, prompt, "")
}

func (p *Provider) generate(ctx context.Context, prompt string, format string) (string, error) {
	// Add timestamp-based seed for variation between requests
	seed := int(time.Now().UnixNano() % 1000000)
	
	reqBody := generateRequest{
		Model:  p.model,
		Prompt: prompt,
		Format: format,
		Stream: false,
		Options: map[string]interface{}{
			"num_ctx":        16384, // 16k context — handles large PDFs + full strict prompt
			"temperature":    0.3,   // Slightly higher for variation (was 0.2)
			"top_p":          0.9,   // Nucleus sampling
			"repeat_penalty": 1.1,   // Reduce repetition
			"num_predict":    4096,  // Max output tokens — enough for many tasks/tickets/decisions
			"seed":           seed,  // Random seed for variation
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

		startTime := time.Now()
		p.log.Info().
			Str("model", p.model).
			Int("attempt", attempt).
			Int("prompt_len", len(prompt)).
			Int("seed", seed).
			Msg("ollama: sending request")

		resp, err := p.client.Do(req)
		cancel()
		
		if err != nil {
			lastErr = fmt.Errorf("ollama: request attempt %d: %w", attempt, err)
			p.log.Warn().Err(lastErr).Dur("duration", time.Since(startTime)).Msg("ollama: request failed")
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * 3 * time.Second)
			}
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
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * 3 * time.Second)
			}
			continue
		}

		var genResp generateResponse
		if err := json.Unmarshal(body, &genResp); err != nil {
			lastErr = fmt.Errorf("ollama: decode response attempt %d: %w", attempt, err)
			continue
		}

		clean := sanitizeOutput(genResp.Response)
		
		// Validate JSON if format was requested
		if format == "json" && !json.Valid([]byte(clean)) {
			lastErr = fmt.Errorf("ollama: invalid JSON in response attempt %d", attempt)
			p.log.Warn().
				Str("raw", truncate(clean, 300)).
				Msg("ollama: invalid json, retrying")
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * 3 * time.Second)
			}
			continue
		}

		duration := time.Since(startTime)
		p.log.Info().
			Dur("duration", duration).
			Int("response_len", len(clean)).
			Msg("ollama: generation successful")
		
		// Keep model alive in memory
		go p.keepAlive()
		
		return clean, nil
	}

	return "", fmt.Errorf("ollama: failed after %d attempts: %w", maxRetries, lastErr)
}

// keepAlive sends a keep-alive request to prevent model unloading
func (p *Provider) keepAlive() {
	url := fmt.Sprintf("%s/api/generate", p.baseURL)
	reqBody := map[string]interface{}{
		"model":      p.model,
		"keep_alive": keepAlive,
	}
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	_, _ = p.client.Do(req)
}

func sanitizeOutput(input string) string {
	s := strings.TrimSpace(input)
	
	// Remove markdown code blocks
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```JSON")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)
	
	// Remove any leading text before the JSON
	if idx := strings.Index(s, "{"); idx > 0 {
		s = s[idx:]
	}
	
	// Remove any trailing text after the JSON
	if idx := strings.LastIndex(s, "}"); idx > 0 && idx < len(s)-1 {
		s = s[:idx+1]
	}
	
	return strings.TrimSpace(s)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
