package handler

import (
	"net/http"

	"github.com/meetext/backend/internal/domain/ai"
	ucai "github.com/meetext/backend/internal/usecase/ai"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/meetext/backend/pkg/response"
	"github.com/meetext/backend/pkg/validator"
	"github.com/rs/zerolog"
)

type AIHandler struct {
	uc                    *ucai.UseCase
	transcriptionProvider ai.TranscriptionProvider
	log                   zerolog.Logger
}

func NewAIHandler(uc *ucai.UseCase, tp ai.TranscriptionProvider, log zerolog.Logger) *AIHandler {
	return &AIHandler{
		uc:                    uc,
		transcriptionProvider: tp,
		log:                   log.With().Str("component", "ai_handler").Logger(),
	}
}

type AnalyzeTextInput struct {
	Text string `json:"text" validate:"required"`
}

// POST /api/v1/ai/analyze-text
// Directly accepts a block of raw text and processes it via the AI service.
func (h *AIHandler) AnalyzeText(w http.ResponseWriter, r *http.Request) {
	var in AnalyzeTextInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	result, err := h.uc.GenerateMeetingAnalysis(r.Context(), in.Text, nil)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to analyze text")
		response.Error(w, apperr.Wrap(err, http.StatusInternalServerError, "AI_ANALYSIS_FAILED", "Failed to generate AI analysis"))
		return
	}

	response.OK(w, result)
}

type UploadInput struct {
	FileURL string `json:"file_url" validate:"required"`
}

// POST /api/v1/ai/upload
// Placeholder endpoint that simulates a file upload triggering transcription and AI analysis.
func (h *AIHandler) Upload(w http.ResponseWriter, r *http.Request) {
	var in UploadInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	// 1. Simulate Transcription
	// TODO: Replace with actual webhook payload reception from n8n.
	transcript, err := h.transcriptionProvider.Transcribe(r.Context(), in.FileURL)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to transcribe file")
		response.Error(w, apperr.Wrap(err, http.StatusInternalServerError, "TRANSCRIPTION_FAILED", "Failed to transcribe audio"))
		return
	}

	// 2. Perform AI Analysis on the transcript
	result, err := h.uc.GenerateMeetingAnalysis(r.Context(), transcript, nil)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to analyze transcript")
		response.Error(w, apperr.Wrap(err, http.StatusInternalServerError, "AI_ANALYSIS_FAILED", "Failed to generate AI analysis"))
		return
	}

	response.OK(w, map[string]interface{}{
		"transcript": transcript,
		"analysis":   result,
	})
}
