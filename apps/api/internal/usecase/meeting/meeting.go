package meeting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/meetext/backend/internal/domain/ai"
	"github.com/meetext/backend/internal/domain/meeting"
	"github.com/meetext/backend/internal/infrastructure/pdf"
	"github.com/meetext/backend/internal/infrastructure/storage"
	ucai "github.com/meetext/backend/internal/usecase/ai"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/meetext/backend/pkg/constants"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var allowedMIMEs = map[string]meeting.UploadType{
	constants.MIMEAudioMPEG: meeting.UploadTypeAudio,
	constants.MIMEAudioWAV:  meeting.UploadTypeAudio,
	constants.MIMEVideoMP4:  meeting.UploadTypeVideo,
	constants.MIMEAppPDF:    meeting.UploadTypePDF,
}

type UploadInput struct {
	WorkspaceID uuid.UUID
	ProjectID   *uuid.UUID
	ClientID    *uuid.UUID
	UploadedBy  uuid.UUID
	Title       string
	FileName    string
	MIMEType    string
	Size        int64
	Reader      io.Reader
}

type UseCase struct {
	repo         meeting.Repository
	storage      storage.Provider
	aiUC         *ucai.UseCase
	pdfExtractor *pdf.Extractor
	log          zerolog.Logger
}

func NewUseCase(repo meeting.Repository, storage storage.Provider, aiUC *ucai.UseCase, pdfExtractor *pdf.Extractor) *UseCase {
	return &UseCase{
		repo:         repo,
		storage:      storage,
		aiUC:         aiUC,
		pdfExtractor: pdfExtractor,
		log:          log.With().Str("component", "meeting_usecase").Logger(),
	}
}

// Upload validates, stores the file, saves the meeting record immediately,
// then kicks off async AI processing with detailed progress tracking.
func (uc *UseCase) Upload(ctx context.Context, in UploadInput) (*meeting.Meeting, error) {
	if in.Size > constants.MaxUploadBytes {
		return nil, apperr.ErrFileTooLarge
	}

	uploadType, ok := allowedMIMEs[in.MIMEType]
	if !ok {
		return nil, apperr.ErrUnsupportedFile
	}

	// Block audio/video uploads - Whisper not integrated yet
	if uploadType == meeting.UploadTypeAudio || uploadType == meeting.UploadTypeVideo {
		return nil, apperr.ErrAudioVideoUnsupported
	}

	fileBytes, err := io.ReadAll(in.Reader)
	if err != nil {
		return nil, fmt.Errorf("meeting: read file: %w", err)
	}

	id := uuid.New()
	key := fmt.Sprintf("workspaces/%s/meetings/%s_%s", in.WorkspaceID, id, in.FileName)

	fileURL, err := uc.storage.Upload(ctx, key, bytes.NewReader(fileBytes), in.Size, in.MIMEType)
	if err != nil {
		return nil, fmt.Errorf("meeting: upload file: %w", err)
	}

	uploadedBy := &in.UploadedBy
	now := time.Now()
	m := &meeting.Meeting{
		ID:                  id,
		WorkspaceID:         in.WorkspaceID,
		ProjectID:           in.ProjectID,
		ClientID:            in.ClientID,
		Title:               in.Title,
		UploadType:          uploadType,
		OriginalFileURL:     fileURL,
		Status:              meeting.StatusProcessing,
		ProcessingStartedAt: &now,
		UploadedBy:          uploadedBy,
		CreatedAt:           now,
	}

	if err := uc.repo.Create(ctx, m); err != nil {
		return nil, fmt.Errorf("meeting: create record: %w", err)
	}

	// Kick off async AI processing with progress tracking
	go uc.processAsync(m, fileBytes, in)

	return m, nil
}

// processAsync runs the full multi-stage AI pipeline in the background with detailed logging.
func (uc *UseCase) processAsync(m *meeting.Meeting, fileBytes []byte, in UploadInput) {
	ctx := context.Background()
	l := uc.log.With().Str("meeting_id", m.ID.String()).Logger()

	defer func() {
		if r := recover(); r != nil {
			l.Error().Interface("panic", r).Msg("meeting: panic in async processing")
			_ = uc.repo.UpdateStatus(ctx, m.ID, meeting.StatusFailed)
		}
	}()

	l.Info().Msg("meeting: starting async AI processing pipeline")

	// Stage 1: Extract PDF text
	l.Info().Msg("meeting: extracting PDF text")
	txt, err := uc.pdfExtractor.Extract(ctx, bytes.NewReader(fileBytes), in.FileName)
	if err != nil {
		l.Error().Err(err).Msg("meeting: pdf extraction failed")
		_ = uc.repo.UpdateStatus(ctx, m.ID, meeting.StatusFailed)
		return
	}

	if txt == "" {
		l.Warn().Msg("meeting: extracted empty text from PDF")
		_ = uc.repo.UpdateStatus(ctx, m.ID, meeting.StatusFailed)
		return
	}

	// Log first 500 chars to verify different PDFs produce different text
	preview := txt
	if len(preview) > 500 {
		preview = preview[:500]
	}
	l.Info().
		Int("text_length", len(txt)).
		Str("preview", preview).
		Msg("meeting: PDF text extracted successfully")

	// Stage 2-6: Multi-stage AI pipeline with progress tracking
	progressCallback := func(progress ucai.ProcessingProgress) {
		l.Info().
			Str("stage", string(progress.Stage)).
			Int("current_chunk", progress.CurrentChunk).
			Int("total_chunks", progress.TotalChunks).
			Str("message", progress.Message).
			Msg("meeting: AI processing progress")
	}

	aiResult, err := uc.aiUC.GenerateMeetingAnalysis(ctx, txt, progressCallback)
	if err != nil {
		l.Error().Err(err).Msg("meeting: AI analysis failed")
		_ = uc.repo.UpdateStatus(ctx, m.ID, meeting.StatusFailed)
		return
	}

	// Log AI result preview to verify different outputs
	l.Info().
		Str("summary_preview", truncateStr(aiResult.Summary, 200)).
		Int("tasks_count", len(aiResult.Tasks)).
		Int("decisions_count", len(aiResult.Decisions)).
		Msg("meeting: AI analysis completed")

	// Stage 7: Update meeting with results
	completed := time.Now()
	m.Transcript = &txt
	m.AISummary = &aiResult.Summary
	m.Status = meeting.StatusCompleted
	m.ProcessingCompletedAt = &completed

	// Serialize the full AI result to JSON so the frontend can display all
	// structured fields (tasks, decisions, risks, etc.) from a single API response.
	if resultJSON, err := json.Marshal(aiResult); err == nil {
		resultStr := string(resultJSON)
		m.AIResultJSON = &resultStr
	} else {
		l.Warn().Err(err).Msg("meeting: failed to serialize ai result json")
	}

	if err := uc.repo.Update(ctx, m); err != nil {
		l.Error().Err(err).Msg("meeting: update record failed")
		_ = uc.repo.UpdateStatus(ctx, m.ID, meeting.StatusFailed)
		return
	}

	// Stage 8: Save structured outputs
	uc.saveStructuredOutputs(ctx, m.ID, in.WorkspaceID, in.ProjectID, in.Title, aiResult, l)

	l.Info().
		Int("tasks", len(aiResult.Tasks)).
		Int("decisions", len(aiResult.Decisions)).
		Int("risks", len(aiResult.Risks)).
		Int("blockers", len(aiResult.Blockers)).
		Msg("meeting: async AI processing completed successfully")
}

func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func (uc *UseCase) saveStructuredOutputs(ctx context.Context, meetingID, workspaceID uuid.UUID, projectID *uuid.UUID, title string, result *ai.AIResult, l zerolog.Logger) {
	// Save tasks
	for _, t := range result.Tasks {
		desc := t.Description
		priority := "medium"
		if t.Priority != nil {
			priority = *t.Priority
		}
		if err := uc.repo.CreateTask(ctx, &meeting.TaskRelation{
			ID:           uuid.New(),
			WorkspaceID:  workspaceID,
			ProjectID:    projectID,
			MeetingID:    &meetingID,
			Title:        t.Title,
			Description:  desc,
			Status:       "todo",
			Priority:     priority,
			AIGenerated:  true,
			AIConfidence: t.ConfidenceScore,
		}); err != nil {
			l.Warn().Err(err).Str("task", t.Title).Msg("meeting: failed to save task")
		}
	}

	// Save decisions
	for _, d := range result.Decisions {
		if err := uc.repo.CreateDecision(ctx, &meeting.DecisionRelation{
			ID:           uuid.New(),
			WorkspaceID:  workspaceID,
			ProjectID:    projectID,
			MeetingID:    meetingID,
			DecisionText: d.Decision,
		}); err != nil {
			l.Warn().Err(err).Msg("meeting: failed to save decision")
		}
	}

	// Save risks
	for _, r := range result.Risks {
		sev := "medium"
		if r.Severity != nil {
			sev = *r.Severity
		}
		if err := uc.repo.CreateBlocker(ctx, &meeting.BlockerRelation{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			ProjectID:   projectID,
			MeetingID:   meetingID,
			BlockerText: r.Risk,
			Severity:    sev,
			Resolved:    false,
		}); err != nil {
			l.Warn().Err(err).Msg("meeting: failed to save risk")
		}
	}

	// Save blockers
	for _, b := range result.Blockers {
		if err := uc.repo.CreateBlocker(ctx, &meeting.BlockerRelation{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			ProjectID:   projectID,
			MeetingID:   meetingID,
			BlockerText: b.Description,
			Severity:    "high",
			Resolved:    false,
		}); err != nil {
			l.Warn().Err(err).Msg("meeting: failed to save blocker")
		}
	}

	// Save document
	if err := uc.repo.CreateDocument(ctx, &meeting.DocumentRelation{
		ID:            uuid.New(),
		WorkspaceID:   workspaceID,
		ProjectID:     projectID,
		MeetingID:     &meetingID,
		Title:         fmt.Sprintf("%s - AI Documentation", title),
		Type:          "meeting_notes",
		Content:       result.ProjectDocumentationMarkdown,
		GeneratedByAI: true,
	}); err != nil {
		l.Warn().Err(err).Msg("meeting: failed to save document")
	}
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*meeting.Meeting, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]*meeting.Meeting, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListByWorkspace(ctx, workspaceID, limit, offset)
}

func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	m, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	_ = uc.storage.Delete(ctx, m.OriginalFileURL)
	return uc.repo.Delete(ctx, id)
}
