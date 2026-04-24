package task

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/task-manager/api/internal/adapter/http/middleware"
	"github.com/task-manager/api/internal/entity"
	"github.com/task-manager/api/internal/port"
	"github.com/task-manager/api/internal/usecase/task"
)

type mockTaskUseCase struct {
	create func(ctx context.Context, userID string, in task.TaskInput) (*entity.Task, error)
	get    func(ctx context.Context, userID, id string) (*entity.Task, error)
	list   func(ctx context.Context, userID string) ([]*entity.Task, error)
	update func(ctx context.Context, userID, id string, in task.TaskUpdateInput) (*entity.Task, error)
	delete func(ctx context.Context, userID, id string, ts time.Time) error
}

func (m *mockTaskUseCase) Create(ctx context.Context, userID string, in task.TaskInput) (*entity.Task, error) {
	if m.create != nil {
		return m.create(ctx, userID, in)
	}
	return &entity.Task{ID: "t1", UserID: userID, Title: in.Title}, nil
}
func (m *mockTaskUseCase) Get(ctx context.Context, userID, id string) (*entity.Task, error) {
	if m.get != nil {
		return m.get(ctx, userID, id)
	}
	return nil, nil
}
func (m *mockTaskUseCase) List(ctx context.Context, userID string) ([]*entity.Task, error) {
	if m.list != nil {
		return m.list(ctx, userID)
	}
	return nil, nil
}
func (m *mockTaskUseCase) Update(ctx context.Context, userID, id string, in task.TaskUpdateInput) (*entity.Task, error) {
	if m.update != nil {
		return m.update(ctx, userID, id, in)
	}
	return nil, nil
}
func (m *mockTaskUseCase) Delete(ctx context.Context, userID, id string, ts time.Time) error {
	if m.delete != nil {
		return m.delete(ctx, userID, id, ts)
	}
	return nil
}

func authCtx(userID, email string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), middleware.UserContextKey(), &middleware.UserInfo{UserID: userID, Email: email})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func TestHandler_Create_Unauthorized(t *testing.T) {
	h := NewHandler(&mockTaskUseCase{})
	r := chi.NewRouter()
	r.Post("/tasks", h.Create)
	body := `{"title":"T","content":"C","due_date":"2025-09-30","request_timestamp":"2025-09-25T20:00:00Z"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandler_Create_Success(t *testing.T) {
	mock := &mockTaskUseCase{
		create: func(ctx context.Context, userID string, in task.TaskInput) (*entity.Task, error) {
			return &entity.Task{ID: "t1", UserID: userID, Title: in.Title, Content: in.Content, DueDate: in.DueDate, RequestTimestamp: in.RequestTimestamp}, nil
		},
	}
	h := NewHandler(mock)
	r := chi.NewRouter()
	r.With(authCtx("u1", "a@b.com")).Post("/tasks", h.Create)
	body := `{"title":"Write","content":"Prepare lesson","due_date":"2025-09-30","request_timestamp":"2025-09-25T20:00:00Z"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var resp TaskResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Equal(t, "t1", resp.ID)
	assert.Equal(t, "Write", resp.Title)
}

func TestHandler_Get_NotFound(t *testing.T) {
	mock := &mockTaskUseCase{
		get: func(ctx context.Context, userID, id string) (*entity.Task, error) { return nil, port.ErrNotFound },
	}
	h := NewHandler(mock)
	r := chi.NewRouter()
	r.Route("/tasks", func(r chi.Router) {
		r.With(authCtx("u1", "a@b.com")).Get("/{id}", h.Get)
	})
	req := httptest.NewRequest("GET", "/tasks/nonexistent", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestHandler_Update_Conflict(t *testing.T) {
	mock := &mockTaskUseCase{
		update: func(ctx context.Context, userID, id string, in task.TaskUpdateInput) (*entity.Task, error) {
			return nil, port.ErrConflict
		},
	}
	h := NewHandler(mock)
	r := chi.NewRouter()
	r.Route("/tasks", func(r chi.Router) {
		r.With(authCtx("u1", "a@b.com")).Put("/{id}", h.Update)
	})
	body := `{"title":"New","request_timestamp":"2025-09-25T20:01:00Z"}`
	req := httptest.NewRequest("PUT", "/tasks/t1", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestHandler_Delete_InvalidBody(t *testing.T) {
	h := NewHandler(&mockTaskUseCase{})
	r := chi.NewRouter()
	r.Route("/tasks", func(r chi.Router) {
		r.With(authCtx("u1", "a@b.com")).Delete("/{id}", h.Delete)
	})
	req := httptest.NewRequest("DELETE", "/tasks/t1", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
