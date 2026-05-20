package constants

type contextKey string

const (
	// Context keys
	CtxUserID      contextKey = "user_id"
	CtxWorkspaceID contextKey = "workspace_id"
	CtxUserRole    contextKey = "user_role"

	// Workspace roles
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"

	// Meeting status
	MeetingStatusPending    = "pending"
	MeetingStatusProcessing = "processing"
	MeetingStatusDone       = "done"
	MeetingStatusFailed     = "failed"

	// File upload limits
	MaxUploadBytes = 2 << 30 // 2 GB

	// Supported MIME types for meeting uploads
	MIMEAudioMPEG = "audio/mpeg"
	MIMEAudioWAV  = "audio/wav"
	MIMEVideoMP4  = "video/mp4"
	MIMEAppPDF    = "application/pdf"
)
