package auth

import (
	"context"
	"net/http"

	apphttp "github.com/task-manager/api/internal/adapter/http"
	"github.com/task-manager/api/internal/entity"
	"github.com/task-manager/api/internal/usecase/auth"
	"github.com/task-manager/api/pkg/jwtauth"
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

func (h *Handler) respondWithToken(w http.ResponseWriter, user *entity.User, status int) bool {
	token, err := h.signer.Sign(user.ID, user.Email)
	if err != nil {
		apphttp.Error(w, http.StatusInternalServerError, "internal server error")
		return false
	}
	apphttp.JSON(w, status, AuthResponse{Token: token})
	return true
}

// Register handles POST /auth/register.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if !apphttp.ReadAndValidateJSON(w, r, &req) {
		return
	}
	user, err := h.authUC.Register(r.Context(), auth.AuthInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if !apphttp.WritePortError(w, err) {
			apphttp.Error(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	h.respondWithToken(w, user, http.StatusCreated)
}

// Login handles POST /auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if !apphttp.ReadAndValidateJSON(w, r, &req) {
		return
	}
	user, err := h.authUC.Login(r.Context(), auth.AuthInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if !apphttp.WritePortError(w, err) {
			apphttp.Error(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	h.respondWithToken(w, user, http.StatusOK)
}
