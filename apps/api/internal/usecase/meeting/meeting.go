package meeting

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/meetext/backend/internal/domain/meeting"
	domainai "github.com/meetext/backend/internal/domain/ai"
	"github.com/meetext/backend/internal/infrastructure/pdf"
	"github.com/meetext/backend/internal/infrastructure/storage"
	ucai "github.com/meetext/backend/internal/usecase/ai"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/meetext/backend/pkg/constants"
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
}

func NewUseCase(repo meeting.Repository, storage storage.Provider, aiUC *ucai.UseCase, pdfExtractor *pdf.Extractor) *UseCase {
	return &UseCase{
		repo:         repo,
		storage:      storage,
		aiUC:         aiUC,
		pdfExtractor: pdfExtractor,
	}
}

func (uc *UseCase) Upload(ctx context.Context, in UploadInput) (*meeting.Meeting, *domainai.AIResult, error) {
	if in.Size > constants.MaxUploadBytes {
		return nil, nil, apperr.ErrFileTooLarge
	}

	uploadType, ok := allowedMIMEs[in.MIMEType]
	if !ok {
		return nil, nil, apperr.ErrUnsupportedFile
	}

	// Reject audio and video uploads gracefully since they are coming soon
	if uploadType == meeting.UploadTypeAudio || uploadType == meeting.UploadTypeVideo {
		return nil, nil, apperr.ErrAudioVideoUnsupported
	}

	id := uuid.New()
	key := fmt.Sprintf("workspaces/%s/meetings/%s_%s", in.WorkspaceID, id, in.FileName)

	// Read all file bytes first so we can use them both for storage and PDF text extraction
	fileBytes, err := io.ReadAll(in.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("meeting: read file bytes: %w", err)
	}

	storageReader := bytes.NewReader(fileBytes)
	fileURL, err := uc.storage.Upload(ctx, key, storageReader, in.Size, in.MIMEType)
	if err != nil {
		return nil, nil, fmt.Errorf("meeting: upload file: %w", err)
	}

	uploadedBy := &in.UploadedBy
	m := &meeting.Meeting{
		ID:              id,
		WorkspaceID:     in.WorkspaceID,
		ProjectID:       in.ProjectID,
		ClientID:        in.ClientID,
		Title:           in.Title,
		UploadType:      uploadType,
		OriginalFileURL: fileURL,
		Status:          meeting.StatusUploaded,
		UploadedBy:      uploadedBy,
		CreatedAt:       time.Now(),
	}

	// For PDF files, extract text and generate analysis via Ollama synchronously
	if uploadType == meeting.UploadTypePDF {
		started := time.Now()
		m.ProcessingStartedAt = &started
		m.Status = meeting.StatusProcessing

		extractorReader := bytes.NewReader(fileBytes)
		txt, err := uc.pdfExtractor.Extract(ctx, extractorReader, in.FileName)
		if err != nil {
			m.Status = meeting.StatusFailed
			_ = uc.repo.Create(ctx, m)
			return nil, nil, fmt.Errorf("meeting: extract pdf text: %w", err)
		}

		m.Transcript = &txt

		aiResult, err := uc.aiUC.GenerateMeetingAnalysis(ctx, txt)
		if err != nil {
			m.Status = meeting.StatusFailed
			_ = uc.repo.Create(ctx, m)
			return nil, nil, fmt.Errorf("meeting: ai meeting analysis: %w", err)
		}

		m.AISummary = &aiResult.Summary
		m.Status = meeting.StatusCompleted
		completed := time.Now()
		m.ProcessingCompletedAt = &completed

		if err := uc.repo.Create(ctx, m); err != nil {
			return nil, nil, fmt.Errorf("meeting: create record: %w", err)
		}

		// Save tasks
		for _, t := range aiResult.Tasks {
			taskID := uuid.New()
			desc := t.Description
			priority := t.Priority
			if priority == "" {
				priority = "medium"
			}
			taskRel := &meeting.TaskRelation{
				ID:           taskID,
				WorkspaceID:  in.WorkspaceID,
				ProjectID:    in.ProjectID,
				MeetingID:    &id,
				Title:        t.Title,
				Description:  desc,
				Status:       "todo",
				Priority:     priority,
				AIGenerated:  true,
				AIConfidence: 92.0,
			}
			_ = uc.repo.CreateTask(ctx, taskRel)
		}

		// Save decisions
		for _, d := range aiResult.Decisions {
			decID := uuid.New()
			decRel := &meeting.DecisionRelation{
				ID:           decID,
				WorkspaceID:  in.WorkspaceID,
				ProjectID:    in.ProjectID,
				MeetingID:    id,
				DecisionText: d.Description,
			}
			_ = uc.repo.CreateDecision(ctx, decRel)
		}

		// Save blockers/risks
		for _, r := range aiResult.Risks {
			blockerID := uuid.New()
			sev := r.Severity
			if sev == "" {
				sev = "medium"
			}
			blockRel := &meeting.BlockerRelation{
				ID:          blockerID,
				WorkspaceID: in.WorkspaceID,
				ProjectID:   in.ProjectID,
				MeetingID:   id,
				BlockerText: r.Description,
				Severity:    sev,
				Resolved:    false,
			}
			_ = uc.repo.CreateBlocker(ctx, blockRel)
		}

		// Save document relation
		docID := uuid.New()
		docTitle := fmt.Sprintf("%s - Sprint Plan", in.Title)
		docRel := &meeting.DocumentRelation{
			ID:            docID,
			WorkspaceID:   in.WorkspaceID,
			ProjectID:     in.ProjectID,
			MeetingID:     &id,
			Title:         docTitle,
			Type:          "sprint_plan",
			Content:       aiResult.ProjectDocumentation,
			GeneratedByAI: true,
		}
		_ = uc.repo.CreateDocument(ctx, docRel)

		return m, aiResult, nil

	} else {
		if err := uc.repo.Create(ctx, m); err != nil {
			return nil, nil, fmt.Errorf("meeting: create record: %w", err)
		}
	}

	return m, nil, nil
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
