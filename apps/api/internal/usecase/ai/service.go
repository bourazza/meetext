package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/meetext/backend/internal/domain/ai"
	"github.com/meetext/backend/internal/infrastructure/ollama/prompts"
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

// GenerateMeetingAnalysis takes a raw text string, formats the LLM prompt, and calls the provider to get structured JSON.
func (uc *UseCase) GenerateMeetingAnalysis(ctx context.Context, text string) (*ai.AIResult, error) {
	if text == "" {
		return nil, fmt.Errorf("input text cannot be empty")
	}

	prompt := prompts.BuildMeetingAnalysisPrompt(text)

	uc.log.Info().Msg("generating AI analysis from text")

	// Call the LLM provider
	jsonResp, err := uc.llmProvider.GenerateJSON(ctx, prompt)
	if err != nil {
		uc.log.Error().Err(err).Msg("failed to generate JSON from LLM")
		return nil, fmt.Errorf("llm generation failed: %w", err)
	}

	uc.log.Debug().Str("raw_json", jsonResp).Msg("received raw JSON from LLM")

	var result ai.AIResult
	if err := json.Unmarshal([]byte(jsonResp), &result); err != nil {
		uc.log.Error().Err(err).Msg("failed to parse structured JSON from LLM")
		// Often LLMs return messy JSON even with format flags.
		// A production system might add a secondary "repair JSON" LLM pass here.
		return nil, fmt.Errorf("failed to parse structured LLM output: %w", err)
	}

	uc.log.Info().Msg("successfully generated and parsed AI analysis")
	return &result, nil
}
