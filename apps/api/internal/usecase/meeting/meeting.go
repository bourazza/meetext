package meeting

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/meetext/backend/internal/domain/meeting"
	"github.com/meetext/backend/internal/infrastructure/storage"
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
	repo    meeting.Repository
	storage storage.Provider
}

func NewUseCase(repo meeting.Repository, storage storage.Provider) *UseCase {
	return &UseCase{repo: repo, storage: storage}
}

func (uc *UseCase) Upload(ctx context.Context, in UploadInput) (*meeting.Meeting, error) {
	if in.Size > constants.MaxUploadBytes {
		return nil, apperr.ErrFileTooLarge
	}

	uploadType, ok := allowedMIMEs[in.MIMEType]
	if !ok {
		return nil, apperr.ErrUnsupportedFile
	}

	id := uuid.New()
	key := fmt.Sprintf("workspaces/%s/meetings/%s_%s", in.WorkspaceID, id, in.FileName)

	fileURL, err := uc.storage.Upload(ctx, key, in.Reader, in.Size, in.MIMEType)
	if err != nil {
		return nil, fmt.Errorf("meeting: upload file: %w", err)
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

	if err := uc.repo.Create(ctx, m); err != nil {
		return nil, fmt.Errorf("meeting: create record: %w", err)
	}

	return m, nil
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
