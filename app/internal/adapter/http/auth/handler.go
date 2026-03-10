package auth

import (
	"context"
	"encoding/json"
	"net/http"

	apphttp "github.com/task-manager/api/internal/adapter/http"
	"github.com/task-manager/api/internal/entity"
	"github.com/task-manager/api/internal/port"
	"github.com/task-manager/api/internal/usecase/auth"
	"github.com/task-manager/api/pkg/jwtauth"
	"github.com/task-manager/api/pkg/validator"
)

// UseCase defines auth operations needed by Handler.
type UseCase interface {
	Register(ctx context.Context, in auth.AuthInput) (*entity.User, error)
	Login(ctx context.Context, in auth.AuthInput) (*entity.User, error)
}

// Handler handles auth endpoints.
type Handler struct {
	authUC UseCase
	signer *jwtauth.Signer
}

// NewHandler creates a new auth Handler.
func NewHandler(authUC UseCase, signer *jwtauth.Signer) *Handler {
	return &Handler{authUC: authUC, signer: signer}
}

// RegisterRequest is the body for POST /auth/register.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest is the body for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse is the response for register/login.
type AuthResponse struct {
	Token string `json:"token"`
}

// Register handles POST /auth/register.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apphttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate(&req); err != nil {
		apphttp.Error(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}
	user, err := h.authUC.Register(r.Context(), auth.AuthInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if err == port.ErrDuplicateUser {
			apphttp.Error(w, http.StatusConflict, "user already exists")
			return
		}
		apphttp.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	token, err := h.signer.Sign(user.ID, user.Email)
	if err != nil {
		apphttp.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	apphttp.JSON(w, http.StatusCreated, AuthResponse{Token: token})
}

// Login handles POST /auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apphttp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate(&req); err != nil {
		apphttp.Error(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}
	user, err := h.authUC.Login(r.Context(), auth.AuthInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if err == port.ErrUnauthorized {
			apphttp.Error(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		apphttp.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	token, err := h.signer.Sign(user.ID, user.Email)
	if err != nil {
		apphttp.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	apphttp.JSON(w, http.StatusOK, AuthResponse{Token: token})
}
