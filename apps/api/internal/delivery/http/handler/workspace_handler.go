package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/meetext/backend/internal/usecase/workspace"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/meetext/backend/pkg/constants"
	"github.com/meetext/backend/pkg/response"
	"github.com/meetext/backend/pkg/validator"
)

type WorkspaceHandler struct {
	uc *workspace.UseCase
}

func NewWorkspaceHandler(uc *workspace.UseCase) *WorkspaceHandler {
	return &WorkspaceHandler{uc: uc}
}

// GET /api/v1/workspaces
func (h *WorkspaceHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(constants.CtxUserID).(uuid.UUID)
	workspaces, err := h.uc.ListForUser(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, workspaces)
}

// GET /api/v1/workspaces/{workspaceID}
func (h *WorkspaceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "workspaceID"))
	if err != nil {
		response.Error(w, apperr.ErrBadRequest)
		return
	}
	ws, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, ws)
}


// PATCH /api/v1/workspaces/{workspaceID}
func (h *WorkspaceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "workspaceID"))
	if err != nil {
		response.Error(w, apperr.ErrBadRequest)
		return
	}

	var body struct {
		Name string `json:"name" validate:"required,min=2,max=100"`
	}
	fields, err := validator.Decode(r, &body)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	userID, _ := r.Context().Value(constants.CtxUserID).(uuid.UUID)
	ws, err := h.uc.UpdateName(r.Context(), id, body.Name, userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, ws)
}

// GET /api/v1/workspaces/{workspaceID}/members
func (h *WorkspaceHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "workspaceID"))
	if err != nil {
		response.Error(w, apperr.ErrBadRequest)
		return
	}
	members, err := h.uc.ListMembers(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, members)
}

// DELETE /api/v1/workspaces/{workspaceID}/members/{userID}
func (h *WorkspaceHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	workspaceID, err := uuid.Parse(chi.URLParam(r, "workspaceID"))
	if err != nil {
		response.Error(w, apperr.ErrBadRequest)
		return
	}
	targetUserID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		response.Error(w, apperr.ErrBadRequest)
		return
	}

	requesterID, _ := r.Context().Value(constants.CtxUserID).(uuid.UUID)
	if err := h.uc.RemoveMember(r.Context(), workspaceID, targetUserID, requesterID); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
