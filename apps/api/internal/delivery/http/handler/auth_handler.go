package handler

import (
	"net/http"

	"github.com/google/uuid"
	ucauth "github.com/meetext/backend/internal/usecase/auth"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/meetext/backend/pkg/constants"
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

// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	response.OK(w, map[string]string{"message": "Signed out."})
}

// POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var in ucauth.ForgotPasswordInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.uc.RequestPasswordReset(r.Context(), in); err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, map[string]string{"message": "If that email exists, a reset link has been sent."})
}

// POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var in ucauth.ResetPasswordInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.uc.ResetPassword(r.Context(), in); err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, map[string]string{"message": "Your password has been updated."})
}

// POST /api/v1/auth/verify-email
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var in ucauth.VerifyEmailInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.uc.VerifyEmail(r.Context(), in); err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, map[string]string{"message": "Email verified."})
}

// POST /api/v1/auth/resend-verification
func (h *AuthHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var in ucauth.ResendVerificationInput
	fields, err := validator.Decode(r, &in)
	if fields != nil {
		response.ValidationError(w, fields)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.uc.ResendVerification(r.Context(), in); err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, map[string]string{"message": "If that account needs verification, a new link has been sent."})
}

// GET /api/v1/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constants.CtxUserID).(uuid.UUID)
	if !ok {
		response.Error(w, apperr.ErrUnauthorized)
		return
	}

	u, err := h.uc.CurrentUser(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, u)
}
