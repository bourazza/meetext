package transcript

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Segment struct {
	Start   float64 `json:"start"`
	End     float64 `json:"end"`
	Speaker *string `json:"speaker,omitempty"`
	Text    string  `json:"text"`
}

type Transcript struct {
	ID         uuid.UUID  `json:"id"`
	MeetingID  uuid.UUID  `json:"meeting_id"`
	RawText    string     `json:"raw_text"`
	Segments   []Segment  `json:"segments"`
	Language   string     `json:"language"`
	ModelUsed  string     `json:"model_used"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, t *Transcript) error
	GetByMeetingID(ctx context.Context, meetingID uuid.UUID) (*Transcript, error)
	Update(ctx context.Context, t *Transcript) error
	Delete(ctx context.Context, meetingID uuid.UUID) error
}
