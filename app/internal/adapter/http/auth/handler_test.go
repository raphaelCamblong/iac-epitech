package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/task-manager/api/internal/entity"
	"github.com/task-manager/api/internal/port"
	"github.com/task-manager/api/internal/usecase/auth"
	"github.com/task-manager/api/pkg/jwtauth"
)

var _ UseCase = (*mockAuthUseCase)(nil)

type mockAuthUseCase struct {
	register func(ctx context.Context, in auth.AuthInput) (*entity.User, error)
	login    func(ctx context.Context, in auth.AuthInput) (*entity.User, error)
}

func (m *mockAuthUseCase) Register(ctx context.Context, in auth.AuthInput) (*entity.User, error) {
	if m.register != nil {
		return m.register(ctx, in)
	}
	return &entity.User{ID: "u1", Email: in.Email}, nil
}
func (m *mockAuthUseCase) Login(ctx context.Context, in auth.AuthInput) (*entity.User, error) {
	if m.login != nil {
		return m.login(ctx, in)
	}
	return nil, nil
}

func TestHandler_Register_Success(t *testing.T) {
	signer := jwtauth.NewSigner("secret", 0)
	mock := &mockAuthUseCase{
		register: func(ctx context.Context, in auth.AuthInput) (*entity.User, error) {
			return &entity.User{ID: "u1", Email: in.Email}, nil
		},
	}
	h := NewHandler(mock, signer)
	r := chi.NewRouter()
	r.Post("/auth/register", h.Register)
	body := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var resp AuthResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.NotEmpty(t, resp.Token)
}

func TestHandler_Register_InvalidBody(t *testing.T) {
	h := NewHandler(&mockAuthUseCase{}, jwtauth.NewSigner("s", 0))
	r := chi.NewRouter()
	r.Post("/auth/register", h.Register)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Register_Duplicate(t *testing.T) {
	mock := &mockAuthUseCase{
		register: func(ctx context.Context, in auth.AuthInput) (*entity.User, error) {
			return nil, port.ErrDuplicateUser
		},
	}
	h := NewHandler(mock, jwtauth.NewSigner("s", 0))
	r := chi.NewRouter()
	r.Post("/auth/register", h.Register)
	body := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestHandler_Login_Unauthorized(t *testing.T) {
	mock := &mockAuthUseCase{
		login: func(ctx context.Context, in auth.AuthInput) (*entity.User, error) {
			return nil, port.ErrUnauthorized
		},
	}
	h := NewHandler(mock, jwtauth.NewSigner("s", 0))
	r := chi.NewRouter()
	r.Post("/auth/login", h.Login)
	body := `{"email":"test@example.com","password":"wrong"}`
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
