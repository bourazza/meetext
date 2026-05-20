package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/meetext/backend/internal/domain/ai"
	"github.com/meetext/backend/internal/infrastructure/ollama/prompts"
	infrapdf "github.com/meetext/backend/internal/infrastructure/pdf"
	"github.com/rs/zerolog"
)

type UseCase struct {
	llmProvider ai.LLMProvider
	log         zerolog.Logger
}

func NewUseCase(llmProvider ai.LLMProvider, log zerolog.Logger) *UseCase {
	return &UseCase{
		llmProvider: llmProvider,
		log:         log.With().Str("component", "ai_usecase").Logger(),
	}
}

// GenerateMeetingAnalysis implements a map-reduce pipeline:
// 1. Chunk text into safe sizes
// 2. Summarize each chunk independently
// 3. Merge summaries
// 4. Run final structured extraction on merged summary
func (uc *UseCase) GenerateMeetingAnalysis(ctx context.Context, text string) (*ai.AIResult, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("ai: input text is empty")
	}

	chunks := infrapdf.ChunkText(text)
	uc.log.Info().Int("chunks", len(chunks)).Int("text_len", len(text)).Msg("ai: starting map-reduce pipeline")

	// Stage 1: Summarize each chunk
	var chunkSummaries []string
	for i, chunk := range chunks {
		uc.log.Info().Int("chunk", i+1).Int("total", len(chunks)).Msg("ai: summarizing chunk")
		summary, err := uc.summarizeChunk(ctx, chunk)
		if err != nil {
			uc.log.Warn().Err(err).Int("chunk", i+1).Msg("ai: chunk summary failed, using raw chunk")
			// Fallback: use a truncated version of the raw chunk
			words := strings.Fields(chunk)
			if len(words) > 500 {
				words = words[:500]
			}
			summary = strings.Join(words, " ")
		}
		chunkSummaries = append(chunkSummaries, summary)
	}

	// Stage 2: Merge summaries
	merged := strings.Join(chunkSummaries, "\n\n---\n\n")
	uc.log.Info().Int("merged_len", len(merged)).Msg("ai: merged chunk summaries")

	// Stage 3: Final structured extraction on merged summary
	result, err := uc.extractStructured(ctx, merged)
	if err != nil {
		return nil, fmt.Errorf("ai: structured extraction failed: %w", err)
	}

	uc.log.Info().Msg("ai: map-reduce pipeline complete")
	return result, nil
}

func (uc *UseCase) summarizeChunk(ctx context.Context, chunk string) (string, error) {
	prompt := prompts.BuildChunkSummaryPrompt(chunk)
	raw, err := uc.llmProvider.GenerateJSON(ctx, prompt)
	if err != nil {
		return "", err
	}
	var resp struct {
		Summary string `json:"summary"`
	}
	if err := json.Unmarshal([]byte(raw), &resp); err != nil || resp.Summary == "" {
		// If JSON parse fails, try to use the raw response as plain text
		return strings.TrimSpace(raw), nil
	}
	return resp.Summary, nil
}

func (uc *UseCase) extractStructured(ctx context.Context, text string) (*ai.AIResult, error) {
	prompt := prompts.BuildMeetingAnalysisPrompt(text)
	raw, err := uc.llmProvider.GenerateJSON(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("llm generation failed: %w", err)
	}

	var result ai.AIResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		uc.log.Error().Err(err).Str("raw", raw[:min(300, len(raw))]).Msg("ai: failed to parse structured output")
		return nil, fmt.Errorf("ai: parse structured output: %w", err)
	}
	return &result, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
