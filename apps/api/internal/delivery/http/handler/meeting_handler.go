package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/meetext/backend/internal/usecase/meeting"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/meetext/backend/pkg/constants"
	"github.com/meetext/backend/pkg/response"
	"github.com/rs/zerolog"
)

type MeetingHandler struct {
	uc  *meeting.UseCase
	log zerolog.Logger
}

func NewMeetingHandler(uc *meeting.UseCase, log zerolog.Logger) *MeetingHandler {
	return &MeetingHandler{uc: uc, log: log.With().Str("component", "meeting_handler").Logger()}
}

// POST /api/v1/workspaces/{workspaceID}/meetings
func (h *MeetingHandler) Upload(w http.ResponseWriter, r *http.Request) {
	workspaceID, err := uuid.Parse(chi.URLParam(r, "workspaceID"))
	if err != nil {
		response.Error(w, apperr.ErrBadRequest)
		return
	}

	userID, _ := r.Context().Value(constants.CtxUserID).(uuid.UUID)

	if err := r.ParseMultipartForm(constants.MaxUploadBytes); err != nil {
		h.log.Error().Err(err).Str("workspace_id", workspaceID.String()).Msg("meeting: parse multipart form failed")
		response.Error(w, apperr.Wrap(err, http.StatusBadRequest, "BAD_REQUEST", "Failed to parse multipart form"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, apperr.Wrap(err, http.StatusBadRequest, "MISSING_FILE", "file field is required"))
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	if title == "" {
		title = header.Filename
	}

	var projectID *uuid.UUID
	if raw := r.FormValue("project_id"); raw != "" {
		if pid, err := uuid.Parse(raw); err == nil {
			projectID = &pid
		}
	}

	var clientID *uuid.UUID
	if raw := r.FormValue("client_id"); raw != "" {
		if cid, err := uuid.Parse(raw); err == nil {
			clientID = &cid
		}
	}

	h.log.Info().
		Str("workspace_id", workspaceID.String()).
		Str("filename", header.Filename).
		Str("mime", header.Header.Get("Content-Type")).
		Int64("size", header.Size).
		Msg("meeting: upload request received")

	m, err := h.uc.Upload(r.Context(), meeting.UploadInput{
		WorkspaceID: workspaceID,
		ProjectID:   projectID,
		ClientID:    clientID,
		UploadedBy:  userID,
		Title:       title,
		FileName:    header.Filename,
		MIMEType:    header.Header.Get("Content-Type"),
		Size:        header.Size,
		Reader:      file,
	})
	if err != nil {
		h.log.Error().Err(err).Str("workspace_id", workspaceID.String()).Msg("meeting: upload failed")
		response.Error(w, err)
		return
	}

	h.log.Info().Str("meeting_id", m.ID.String()).Str("status", string(m.Status)).Msg("meeting: upload accepted, processing async")
	response.Created(w, map[string]interface{}{
		"meeting": m,
		"message": "Meeting uploaded. AI processing has started in the background.",
	})
}

// GET /api/v1/workspaces/{workspaceID}/meetings
func (h *MeetingHandler) List(w http.ResponseWriter, r *http.Request) {
	workspaceID, err := uuid.Parse(chi.URLParam(r, "workspaceID"))
	if err != nil {
		response.Error(w, apperr.ErrBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	meetings, err := h.uc.List(r.Context(), workspaceID, limit, offset)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, meetings)
}

// GET /api/v1/workspaces/{workspaceID}/meetings/{meetingID}
func (h *MeetingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	meetingID, err := uuid.Parse(chi.URLParam(r, "meetingID"))
	if err != nil {
		response.Error(w, apperr.ErrBadRequest)
		return
	}

	m, err := h.uc.GetByID(r.Context(), meetingID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, m)
}

// GET /api/v1/workspaces/{workspaceID}/meetings/{meetingID}/status
// Lightweight polling endpoint — returns only id + status + summary.
// Validates workspace ownership to prevent unauthorized access.
func (h *MeetingHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	workspaceID, err := uuid.Parse(chi.URLParam(r, "workspaceID"))
	if err != nil {
		h.log.Warn().Str("workspace_id", chi.URLParam(r, "workspaceID")).Msg("meeting: invalid workspace ID")
		response.Error(w, apperr.ErrBadRequest)
		return
	}

	meetingID, err := uuid.Parse(chi.URLParam(r, "meetingID"))
	if err != nil {
		h.log.Warn().Str("meeting_id", chi.URLParam(r, "meetingID")).Msg("meeting: invalid meeting ID")
		response.Error(w, apperr.ErrBadRequest)
		return
	}

	userID, _ := r.Context().Value(constants.CtxUserID).(uuid.UUID)

	h.log.Debug().
		Str("workspace_id", workspaceID.String()).
		Str("meeting_id", meetingID.String()).
		Str("user_id", userID.String()).
		Msg("meeting: status poll request")

	m, err := h.uc.GetByID(r.Context(), meetingID)
	if err != nil {
		h.log.Warn().Err(err).Str("meeting_id", meetingID.String()).Msg("meeting: status fetch failed")
		response.Error(w, err)
		return
	}

	// Validate workspace ownership
	if m.WorkspaceID != workspaceID {
		h.log.Warn().
			Str("meeting_workspace", m.WorkspaceID.String()).
			Str("requested_workspace", workspaceID.String()).
			Msg("meeting: workspace mismatch")
		response.Error(w, apperr.ErrForbidden)
		return
	}

	response.OK(w, map[string]interface{}{
		"id":         m.ID,
		"status":     m.Status,
		"ai_summary": m.AISummary,
		"ai_result":  m.AIResultJSON, // Full structured JSON: tasks, decisions, risks, etc.
	})
}

// DELETE /api/v1/workspaces/{workspaceID}/meetings/{meetingID}
func (h *MeetingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	meetingID, err := uuid.Parse(chi.URLParam(r, "meetingID"))
	if err != nil {
		response.Error(w, apperr.ErrBadRequest)
		return
	}

	if err := h.uc.Delete(r.Context(), meetingID); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
