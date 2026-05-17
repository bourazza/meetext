package handler

import (
	"net/http"

	ucauth "github.com/meetext/backend/internal/usecase/auth"
	"github.com/meetext/backend/pkg/response"
	"github.com/meetext/backend/pkg/validator"
)

type AuthHandler struct {
	uc *ucauth.UseCase
}

func NewAuthHandler(uc *ucauth.UseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in ucauth.RegisterInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	res, err := h.uc.Register(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, res)
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in ucauth.LoginInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	res, err := h.uc.Login(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, res)
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
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

	tokens, err := h.uc.RefreshToken(r.Context(), body.RefreshToken)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, tokens)
}
