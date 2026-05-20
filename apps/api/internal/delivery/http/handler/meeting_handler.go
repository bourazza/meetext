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
	"github.com/rs/zerolog/log"
)

type MeetingHandler struct {
	uc *meeting.UseCase
}

func NewMeetingHandler(uc *meeting.UseCase) *MeetingHandler {
	return &MeetingHandler{uc: uc}
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
		response.Error(w, apperr.ErrFileTooLarge)
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
		pid, err := uuid.Parse(raw)
		if err == nil {
			projectID = &pid
		}
	}

	var clientID *uuid.UUID
	if raw := r.FormValue("client_id"); raw != "" {
		cid, err := uuid.Parse(raw)
		if err == nil {
			clientID = &cid
		}
	}

	m, aiRes, err := h.uc.Upload(r.Context(), meeting.UploadInput{
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
		log.Error().Err(err).Msg("meeting: upload failed")
		response.Error(w, err)
		return
	}

	if aiRes != nil {
		response.Created(w, map[string]interface{}{
			"meeting":  m,
			"analysis": aiRes,
		})
	} else {
		response.Created(w, map[string]interface{}{
			"meeting": m,
		})
	}
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
