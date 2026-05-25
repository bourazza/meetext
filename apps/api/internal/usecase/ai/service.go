package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/meetext/backend/internal/domain/ai"
	"github.com/meetext/backend/internal/infrastructure/ollama/prompts"
	infrapdf "github.com/meetext/backend/internal/infrastructure/pdf"
	"github.com/rs/zerolog"
)

const (
	maxParallelWorkers = 2 // Limit for Ryzen 5 4800U with 16GB RAM
)

type ProcessingStage string

const (
	StageExtracting  ProcessingStage = "extracting"
	StageCleaning    ProcessingStage = "cleaning"
	StageChunking    ProcessingStage = "chunking"
	StageSummarizing ProcessingStage = "summarizing"
	StageExtracting2 ProcessingStage = "extracting_structured"
	StageGenerating  ProcessingStage = "generating_report"
	StageCompleted   ProcessingStage = "completed"
	StageFailed      ProcessingStage = "failed"
)

type ProcessingProgress struct {
	Stage          ProcessingStage
	CurrentChunk   int
	TotalChunks    int
	Message        string
	CompletedAt    *time.Time
	FailureReason  string
}

type ProgressCallback func(progress ProcessingProgress)

type UseCase struct {
	llmProvider ai.LLMProvider
	log         zerolog.Logger
	promptMode  prompts.PromptMode
}

func NewUseCase(llmProvider ai.LLMProvider, log zerolog.Logger) *UseCase {
	return &UseCase{
		llmProvider: llmProvider,
		log:         log.With().Str("component", "ai_usecase").Logger(),
		promptMode:  prompts.PromptModeStrict, // default to strict: exhaustive extraction, no hallucination
	}
}

// SetPromptMode allows changing the prompt mode (fast, balanced, strict)
func (uc *UseCase) SetPromptMode(mode prompts.PromptMode) {
	uc.promptMode = mode
	uc.log.Info().Str("mode", string(mode)).Msg("ai: prompt mode changed")
}

// GenerateMeetingAnalysis implements a robust multi-stage pipeline:
// 1. Clean text
// 2. Chunk into semantic pieces
// 3. Parallel summarize chunks (with worker pool)
// 4. Merge summaries
// 5. Extract structured data from merged summaries only
// 6. Generate final report
func (uc *UseCase) GenerateMeetingAnalysis(ctx context.Context, rawText string, progressCb ProgressCallback) (*ai.AIResult, error) {
	if strings.TrimSpace(rawText) == "" {
		return nil, fmt.Errorf("ai: input text is empty")
	}

	startTime := time.Now()
	uc.log.Info().Int("raw_text_len", len(rawText)).Msg("ai: starting multi-stage pipeline")

	// Stage 1: Clean text
	if progressCb != nil {
		progressCb(ProcessingProgress{Stage: StageCleaning, Message: "Cleaning transcript"})
	}
	cleanedText := infrapdf.CleanText(rawText)
	uc.log.Info().
		Int("original_len", len(rawText)).
		Int("cleaned_len", len(cleanedText)).
		Msg("ai: text cleaned")

	// Stage 2: Chunk text
	if progressCb != nil {
		progressCb(ProcessingProgress{Stage: StageChunking, Message: "Splitting into semantic chunks"})
	}
	chunks := infrapdf.ChunkText(cleanedText)
	uc.log.Info().
		Int("num_chunks", len(chunks)).
		Msg("ai: text chunked")

	if len(chunks) == 0 {
		return nil, fmt.Errorf("ai: no chunks created from text")
	}

	// Stage 3: Parallel chunk summarization with worker pool
	if progressCb != nil {
		progressCb(ProcessingProgress{
			Stage:        StageSummarizing,
			TotalChunks:  len(chunks),
			CurrentChunk: 0,
			Message:      fmt.Sprintf("Summarizing %d chunks", len(chunks)),
		})
	}

	chunkSummaries, err := uc.summarizeChunksParallel(ctx, chunks, progressCb)
	if err != nil {
		return nil, fmt.Errorf("ai: chunk summarization failed: %w", err)
	}

	// Stage 4: Merge summaries
	merged := strings.Join(chunkSummaries, "\n\n---\n\n")
	uc.log.Info().
		Int("merged_len", len(merged)).
		Int("num_summaries", len(chunkSummaries)).
		Msg("ai: summaries merged")

	// Stage 5: Extract structured data from merged summaries ONLY
	if progressCb != nil {
		progressCb(ProcessingProgress{Stage: StageExtracting2, Message: "Extracting structured data"})
	}
	result, err := uc.extractStructured(ctx, merged)
	if err != nil {
		return nil, fmt.Errorf("ai: structured extraction failed: %w", err)
	}

	// Stage 6: Generate final report
	if progressCb != nil {
		progressCb(ProcessingProgress{Stage: StageGenerating, Message: "Generating final report"})
	}

	duration := time.Since(startTime)
	uc.log.Info().
		Dur("total_duration", duration).
		Int("num_tasks", len(result.Tasks)).
		Int("num_decisions", len(result.Decisions)).
		Int("num_risks", len(result.Risks)).
		Msg("ai: pipeline completed successfully")

	if progressCb != nil {
		completed := time.Now()
		progressCb(ProcessingProgress{
			Stage:       StageCompleted,
			Message:     "Processing complete",
			CompletedAt: &completed,
		})
	}

	return result, nil
}

// summarizeChunksParallel processes chunks in parallel with a worker pool
func (uc *UseCase) summarizeChunksParallel(ctx context.Context, chunks []infrapdf.Chunk, progressCb ProgressCallback) ([]string, error) {
	type chunkResult struct {
		index   int
		summary string
		err     error
	}

	results := make([]chunkResult, len(chunks))
	jobs := make(chan infrapdf.Chunk, len(chunks))
	resultsCh := make(chan chunkResult, len(chunks))

	// Start worker pool
	var wg sync.WaitGroup
	for w := 0; w < maxParallelWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for chunk := range jobs {
				uc.log.Info().
					Int("worker", workerID).
					Int("chunk", chunk.Index+1).
					Int("words", chunk.Words).
					Msg("ai: processing chunk")

				summary, err := uc.summarizeChunk(ctx, chunk.Content)
				resultsCh <- chunkResult{
					index:   chunk.Index,
					summary: summary,
					err:     err,
				}

				if progressCb != nil {
					progressCb(ProcessingProgress{
						Stage:        StageSummarizing,
						CurrentChunk: chunk.Index + 1,
						TotalChunks:  len(chunks),
						Message:      fmt.Sprintf("Processing chunk %d/%d", chunk.Index+1, len(chunks)),
					})
				}
			}
		}(w)
	}

	// Send jobs
	for _, chunk := range chunks {
		jobs <- chunk
	}
	close(jobs)

	// Collect results
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for result := range resultsCh {
		results[result.index] = result
	}

	// Check for errors and build summaries in order
	var summaries []string
	var failedChunks []int
	for i, result := range results {
		if result.err != nil {
			uc.log.Warn().
				Err(result.err).
				Int("chunk", i).
				Msg("ai: chunk summarization failed, using fallback")
			failedChunks = append(failedChunks, i)
			// Fallback: use truncated raw chunk
			words := strings.Fields(chunks[i].Content)
			if len(words) > 300 {
				words = words[:300]
			}
			summaries = append(summaries, strings.Join(words, " ")+"...")
		} else {
			summaries = append(summaries, result.summary)
		}
	}

	if len(failedChunks) > 0 {
		uc.log.Warn().
			Ints("failed_chunks", failedChunks).
			Msg("ai: some chunks failed summarization, using fallback")
	}

	return summaries, nil
}

func (uc *UseCase) summarizeChunk(ctx context.Context, chunk string) (string, error) {
	prompt := prompts.BuildChunkSummaryPrompt(chunk, uc.promptMode)
	
	// Log chunk preview
	uc.log.Debug().
		Str("chunk_preview", truncate(chunk, 200)).
		Msg("ai: summarizing chunk")
	
	raw, err := uc.llmProvider.GenerateJSON(ctx, prompt)
	if err != nil {
		return "", err
	}

	var resp struct {
		Summary string `json:"summary"`
	}
	if err := json.Unmarshal([]byte(raw), &resp); err != nil || resp.Summary == "" {
		// Fallback: use raw response as plain text
		return strings.TrimSpace(raw), nil
	}
	return resp.Summary, nil
}

func (uc *UseCase) extractStructured(ctx context.Context, mergedSummaries string) (*ai.AIResult, error) {
	prompt := prompts.BuildStructuredExtractionPrompt(mergedSummaries, uc.promptMode)
	
	// Log what we're sending to Ollama
	uc.log.Info().
		Int("prompt_len", len(prompt)).
		Int("summaries_len", len(mergedSummaries)).
		Str("summaries_preview", truncate(mergedSummaries, 300)).
		Msg("ai: sending to ollama for structured extraction")
	
	raw, err := uc.llmProvider.GenerateJSON(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("llm generation failed: %w", err)
	}

	// Log what Ollama returned
	uc.log.Info().
		Int("response_len", len(raw)).
		Str("response_preview", truncate(raw, 300)).
		Msg("ai: received from ollama")

	var result ai.AIResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		uc.log.Error().
			Err(err).
			Str("raw", truncate(raw, 500)).
			Msg("ai: failed to parse structured output")
		return nil, fmt.Errorf("ai: parse structured output: %w", err)
	}

	return &result, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
