package interfaces

import "context"

type TranscriptionProvider interface {
    // Transcribe receives a file URL and returns the transcript text.
    // In the MVP, this is only used for future audio/video processing.
    // Implementations may call external services such as n8n + Whisper.
    Transcribe(ctx context.Context, fileURL string) (string, error)
}
