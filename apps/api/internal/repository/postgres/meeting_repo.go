package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meetext/backend/internal/domain/meeting"
	"github.com/meetext/backend/pkg/apperr"
)

type MeetingRepository struct {
	db *pgxpool.Pool
}

func NewMeetingRepository(db *pgxpool.Pool) *MeetingRepository {
	return &MeetingRepository{db: db}
}

func (r *MeetingRepository) Create(ctx context.Context, m *meeting.Meeting) error {
	q := `INSERT INTO meetings
		  (id, workspace_id, project_id, client_id, title, upload_type,
		   original_file_url, status, uploaded_by, created_at)
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.Exec(ctx, q,
		m.ID, m.WorkspaceID, m.ProjectID, m.ClientID, m.Title,
		m.UploadType, m.OriginalFileURL, m.Status, m.UploadedBy, m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("meeting repo: create: %w", err)
	}
	return nil
}

func (r *MeetingRepository) GetByID(ctx context.Context, id uuid.UUID) (*meeting.Meeting, error) {
	q := `SELECT id, workspace_id, project_id, client_id, title, upload_type,
		         original_file_url, transcript, ai_summary, ai_result_json, duration_seconds,
		         language, status, processing_started_at, processing_completed_at,
		         uploaded_by, created_at
		  FROM meetings WHERE id=$1`
	m := &meeting.Meeting{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&m.ID, &m.WorkspaceID, &m.ProjectID, &m.ClientID, &m.Title, &m.UploadType,
		&m.OriginalFileURL, &m.Transcript, &m.AISummary, &m.AIResultJSON, &m.DurationSeconds,
		&m.Language, &m.Status, &m.ProcessingStartedAt, &m.ProcessingCompletedAt,
		&m.UploadedBy, &m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("meeting repo: get by id: %w", err)
	}
	return m, nil
}

func (r *MeetingRepository) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]*meeting.Meeting, error) {
	q := `SELECT id, workspace_id, project_id, client_id, title, upload_type,
		         original_file_url, transcript, ai_summary, ai_result_json, duration_seconds,
		         language, status, processing_started_at, processing_completed_at,
		         uploaded_by, created_at
		  FROM meetings WHERE workspace_id=$1
		  ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	return r.scanMeetings(ctx, q, workspaceID, limit, offset)
}

func (r *MeetingRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*meeting.Meeting, error) {
	q := `SELECT id, workspace_id, project_id, client_id, title, upload_type,
		         original_file_url, transcript, ai_summary, ai_result_json, duration_seconds,
		         language, status, processing_started_at, processing_completed_at,
		         uploaded_by, created_at
		  FROM meetings WHERE project_id=$1 ORDER BY created_at DESC`
	return r.scanMeetings(ctx, q, projectID)
}

func (r *MeetingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status meeting.Status) error {
	q := `UPDATE meetings SET status=$1 WHERE id=$2`
	tag, err := r.db.Exec(ctx, q, status, id)
	if err != nil {
		return fmt.Errorf("meeting repo: update status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *MeetingRepository) Update(ctx context.Context, m *meeting.Meeting) error {
	q := `UPDATE meetings SET title=$1, transcript=$2, ai_summary=$3, ai_result_json=$4,
		  duration_seconds=$5, language=$6, processing_started_at=$7,
		  processing_completed_at=$8, status=$9 WHERE id=$10`
	tag, err := r.db.Exec(ctx, q,
		m.Title, m.Transcript, m.AISummary, m.AIResultJSON, m.DurationSeconds,
		m.Language, m.ProcessingStartedAt, m.ProcessingCompletedAt,
		m.Status, m.ID,
	)
	if err != nil {
		return fmt.Errorf("meeting repo: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *MeetingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM meetings WHERE id=$1`
	tag, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("meeting repo: delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *MeetingRepository) AddParticipant(ctx context.Context, p *meeting.Participant) error {
	q := `INSERT INTO meeting_participants (id, meeting_id, participant_name, participant_email, created_at)
		  VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.Exec(ctx, q, p.ID, p.MeetingID, p.ParticipantName, p.ParticipantEmail, p.CreatedAt)
	if err != nil {
		return fmt.Errorf("meeting repo: add participant: %w", err)
	}
	return nil
}

func (r *MeetingRepository) ListParticipants(ctx context.Context, meetingID uuid.UUID) ([]*meeting.Participant, error) {
	q := `SELECT id, meeting_id, participant_name, participant_email, created_at
		  FROM meeting_participants WHERE meeting_id=$1`
	rows, err := r.db.Query(ctx, q, meetingID)
	if err != nil {
		return nil, fmt.Errorf("meeting repo: list participants: %w", err)
	}
	defer rows.Close()

	var result []*meeting.Participant
	for rows.Next() {
		p := &meeting.Participant{}
		if err := rows.Scan(&p.ID, &p.MeetingID, &p.ParticipantName, &p.ParticipantEmail, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("meeting repo: scan participant: %w", err)
		}
		result = append(result, p)
	}
	return result, nil
}

// scanMeetings is a shared row scanner for meeting queries.
func (r *MeetingRepository) scanMeetings(ctx context.Context, q string, args ...interface{}) ([]*meeting.Meeting, error) {
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("meeting repo: query: %w", err)
	}
	defer rows.Close()

	var result []*meeting.Meeting
	for rows.Next() {
		m := &meeting.Meeting{}
		if err := rows.Scan(
			&m.ID, &m.WorkspaceID, &m.ProjectID, &m.ClientID, &m.Title, &m.UploadType,
			&m.OriginalFileURL, &m.Transcript, &m.AISummary, &m.AIResultJSON, &m.DurationSeconds,
			&m.Language, &m.Status, &m.ProcessingStartedAt, &m.ProcessingCompletedAt,
			&m.UploadedBy, &m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("meeting repo: scan: %w", err)
		}
		result = append(result, m)
	}
	return result, nil
}

func (r *MeetingRepository) CreateTask(ctx context.Context, t *meeting.TaskRelation) error {
	q := `INSERT INTO tasks (id, workspace_id, project_id, meeting_id, title, description, status, priority, due_date, ai_generated, ai_confidence, assigned_to, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())`
	_, err := r.db.Exec(ctx, q, t.ID, t.WorkspaceID, t.ProjectID, t.MeetingID, t.Title, t.Description, t.Status, t.Priority, t.DueDate, t.AIGenerated, t.AIConfidence, t.AssignedTo)
	if err != nil {
		return fmt.Errorf("meeting repo: create task: %w", err)
	}
	return nil
}

func (r *MeetingRepository) CreateDecision(ctx context.Context, d *meeting.DecisionRelation) error {
	q := `INSERT INTO decisions (id, workspace_id, project_id, meeting_id, decision_text, created_at)
		  VALUES ($1, $2, $3, $4, $5, NOW())`
	_, err := r.db.Exec(ctx, q, d.ID, d.WorkspaceID, d.ProjectID, d.MeetingID, d.DecisionText)
	if err != nil {
		return fmt.Errorf("meeting repo: create decision: %w", err)
	}
	return nil
}

func (r *MeetingRepository) CreateBlocker(ctx context.Context, b *meeting.BlockerRelation) error {
	q := `INSERT INTO blockers (id, workspace_id, project_id, meeting_id, blocker_text, severity, resolved, created_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`
	_, err := r.db.Exec(ctx, q, b.ID, b.WorkspaceID, b.ProjectID, b.MeetingID, b.BlockerText, b.Severity, b.Resolved)
	if err != nil {
		return fmt.Errorf("meeting repo: create blocker: %w", err)
	}
	return nil
}

func (r *MeetingRepository) CreateDocument(ctx context.Context, d *meeting.DocumentRelation) error {
	q := `INSERT INTO documents (id, workspace_id, project_id, meeting_id, title, type, content, generated_by_ai, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`
	_, err := r.db.Exec(ctx, q, d.ID, d.WorkspaceID, d.ProjectID, d.MeetingID, d.Title, d.Type, d.Content, d.GeneratedByAI)
	if err != nil {
		return fmt.Errorf("meeting repo: create document: %w", err)
	}
	return nil
}

