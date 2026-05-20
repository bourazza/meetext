package meeting

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
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
// then kicks off async AI processing. Returns the meeting instantly.
func (uc *UseCase) Upload(ctx context.Context, in UploadInput) (*meeting.Meeting, error) {
	if in.Size > constants.MaxUploadBytes {
		return nil, apperr.ErrFileTooLarge
	}

	uploadType, ok := allowedMIMEs[in.MIMEType]
	if !ok {
		return nil, apperr.ErrUnsupportedFile
	}

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

	// Kick off async AI processing — do not block the HTTP response
	go uc.processAsync(m, fileBytes, in)

	return m, nil
}

// processAsync runs the full AI pipeline in the background.
func (uc *UseCase) processAsync(m *meeting.Meeting, fileBytes []byte, in UploadInput) {
	ctx := context.Background()
	l := uc.log.With().Str("meeting_id", m.ID.String()).Logger()

	defer func() {
		if r := recover(); r != nil {
			l.Error().Interface("panic", r).Msg("meeting: panic in async processing")
			_ = uc.repo.UpdateStatus(ctx, m.ID, meeting.StatusFailed)
		}
	}()

	l.Info().Msg("meeting: starting async AI processing")

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

	aiResult, err := uc.aiUC.GenerateMeetingAnalysis(ctx, txt)
	if err != nil {
		l.Error().Err(err).Msg("meeting: AI analysis failed")
		_ = uc.repo.UpdateStatus(ctx, m.ID, meeting.StatusFailed)
		return
	}

	// Update meeting with results
	completed := time.Now()
	m.Transcript = &txt
	m.AISummary = &aiResult.Summary
	m.Status = meeting.StatusCompleted
	m.ProcessingCompletedAt = &completed

	if err := uc.repo.Update(ctx, m); err != nil {
		l.Error().Err(err).Msg("meeting: update record failed")
		_ = uc.repo.UpdateStatus(ctx, m.ID, meeting.StatusFailed)
		return
	}

	// Save tasks
	for _, t := range aiResult.Tasks {
		desc := t.Description
		priority := "medium"
		if t.Priority != nil {
			priority = *t.Priority
		}
		_ = uc.repo.CreateTask(ctx, &meeting.TaskRelation{
			ID:           uuid.New(),
			WorkspaceID:  in.WorkspaceID,
			ProjectID:    in.ProjectID,
			MeetingID:    &m.ID,
			Title:        t.Title,
			Description:  desc,
			Status:       "todo",
			Priority:     priority,
			AIGenerated:  true,
			AIConfidence: t.ConfidenceScore,
		})
	}

	// Save decisions
	for _, d := range aiResult.Decisions {
		_ = uc.repo.CreateDecision(ctx, &meeting.DecisionRelation{
			ID:           uuid.New(),
			WorkspaceID:  in.WorkspaceID,
			ProjectID:    in.ProjectID,
			MeetingID:    m.ID,
			DecisionText: d.Decision,
		})
	}

	// Save blockers/risks
	for _, r := range aiResult.Risks {
		sev := "medium"
		if r.Severity != nil {
			sev = *r.Severity
		}
		_ = uc.repo.CreateBlocker(ctx, &meeting.BlockerRelation{
			ID:          uuid.New(),
			WorkspaceID: in.WorkspaceID,
			ProjectID:   in.ProjectID,
			MeetingID:   m.ID,
			BlockerText: r.Risk,
			Severity:    sev,
			Resolved:    false,
		})
	}

	// Save document
	_ = uc.repo.CreateDocument(ctx, &meeting.DocumentRelation{
		ID:            uuid.New(),
		WorkspaceID:   in.WorkspaceID,
		ProjectID:     in.ProjectID,
		MeetingID:     &m.ID,
		Title:         fmt.Sprintf("%s - AI Documentation", in.Title),
		Type:          "sprint_plan",
		Content:       aiResult.ProjectDocumentationMarkdown,
		GeneratedByAI: true,
	})

	l.Info().Msg("meeting: async AI processing complete")
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
