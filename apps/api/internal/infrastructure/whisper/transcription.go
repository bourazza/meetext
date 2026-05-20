package whisper

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type MockProvider struct {
	log zerolog.Logger
}

func NewMockProvider(log zerolog.Logger) *MockProvider {
	return &MockProvider{
		log: log.With().Str("component", "whisper_mock").Logger(),
	}
}

// Transcribe simulates the transcription of an audio/video file.
// TODO: Replace this mock implementation with the actual n8n + Whisper integration pipeline.
func (p *MockProvider) Transcribe(ctx context.Context, fileURL string) (string, error) {
	p.log.Info().Str("file_url", fileURL).Msg("Simulating transcription for file...")

	// Simulate processing time
	select {
	case <-time.After(2 * time.Second):
		// Simulated transcription
		simulatedText := `
			Okay, let's get started. Thanks everyone for joining.
			First on the agenda is the new API integration for the mobile app. Sarah, you mentioned you'd take the lead on building the authentication endpoint. Can we have that by Friday? Yes, high priority.
			Also, we decided to use PostgreSQL instead of MongoDB for the new microservice. Team agreement on that.
			One risk I want to highlight is that the third-party payment gateway documentation is outdated, which might cause delays. We need to reach out to their support as a mitigation.
			The client, Acme Corp, specifically requested we add an export to PDF feature. We committed to doing it, so let's track that as a feature ticket.
			Finally, David, please follow up with the design team regarding the new logo assets by tomorrow.
		`
		p.log.Info().Msg("Transcription simulation completed.")
		return simulatedText, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
