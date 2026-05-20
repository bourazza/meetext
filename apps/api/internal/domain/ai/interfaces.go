package ai

import "context"

// LLMProvider defines the interface for communicating with an LLM (e.g., Ollama, OpenAI).
type LLMProvider interface {
	// GenerateJSON sends a prompt to the LLM and strictly expects a JSON formatted string in return.
	GenerateJSON(ctx context.Context, prompt string) (string, error)
}

// TranscriptionProvider defines the interface for transcribing audio/video files.
type TranscriptionProvider interface {
	// Transcribe processes the given file and returns the raw transcript text.
	Transcribe(ctx context.Context, fileURL string) (string, error)
}

// Service defines the business logic for AI processing.
type Service interface {
	// GenerateMeetingAnalysis takes a raw transcript or meeting text and generates a structured AIResult.
	GenerateMeetingAnalysis(ctx context.Context, text string) (*AIResult, error)
}
