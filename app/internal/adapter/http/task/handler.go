package task

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	appHttp "github.com/task-manager/api/internal/adapter/http"
	"github.com/task-manager/api/internal/adapter/http/middleware"
	"github.com/task-manager/api/internal/entity"
	"github.com/task-manager/api/internal/usecase/task"
)

// UseCase defines the task operations needed by Handler.
type UseCase interface {
	Create(ctx context.Context, userID string, in task.TaskInput) (*entity.Task, error)
	Get(ctx context.Context, userID, id string) (*entity.Task, error)
	List(ctx context.Context, userID string) ([]*entity.Task, error)
	Update(ctx context.Context, userID, id string, in task.TaskUpdateInput) (*entity.Task, error)
	Delete(ctx context.Context, userID, id string, requestTimestamp time.Time) error
}

// Handler handles task endpoints.
type Handler struct {
	taskUC UseCase
}

// NewHandler creates a new task Handler.
func NewHandler(taskUC UseCase) *Handler {
	return &Handler{taskUC: taskUC}
}

// CreateTaskRequest is the body for POST /tasks.
type CreateTaskRequest struct {
	Title            string `json:"title" validate:"required"`
	Content          string `json:"content" validate:"required"`
	DueDate          string `json:"due_date" validate:"required"`
	RequestTimestamp string `json:"request_timestamp" validate:"required"`
}

// UpdateTaskRequest is the body for PUT /tasks/{id}.
type UpdateTaskRequest struct {
	Title            *string `json:"title"`
	Content          *string `json:"content"`
	DueDate          *string `json:"due_date"`
	Done             *bool   `json:"done"`
	RequestTimestamp string  `json:"request_timestamp" validate:"required"`
}

// DeleteTaskRequest is the body for DELETE /tasks/{id}.
type DeleteTaskRequest struct {
	RequestTimestamp string `json:"request_timestamp" validate:"required"`
}

// TaskResponse represents a task in API responses.
type TaskResponse struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	DueDate          string    `json:"due_date"`
	Done             bool      `json:"done"`
	RequestTimestamp time.Time `json:"request_timestamp"`
}

func toTaskResponse(t *entity.Task) TaskResponse {
	return TaskResponse{
		ID:               t.ID,
		Title:            t.Title,
		Content:          t.Content,
		DueDate:          t.DueDate.Format("2006-01-02"),
		Done:             t.Done,
		RequestTimestamp: t.RequestTimestamp,
	}
}

func userOr401(w http.ResponseWriter, r *http.Request) (*middleware.UserInfo, bool) {
	u := middleware.UserFromContext(r.Context())
	if u == nil {
		appHttp.Error(w, http.StatusUnauthorized, "unauthorized")
		return nil, false
	}
	return u, true
}

func parseDueDate(w http.ResponseWriter, s string) (time.Time, bool) {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		appHttp.Error(w, http.StatusBadRequest, "invalid due_date format")
		return time.Time{}, false
	}
	return d, true
}

func parseRequestTimestamp(w http.ResponseWriter, s string) (time.Time, bool) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		appHttp.Error(w, http.StatusBadRequest, "invalid request_timestamp format")
		return time.Time{}, false
	}
	return t, true
}

// Create handles POST /tasks.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := userOr401(w, r)
	if !ok {
		return
	}
	var req CreateTaskRequest
	if !appHttp.ReadAndValidateJSON(w, r, &req) {
		return
	}
	dueDate, ok := parseDueDate(w, req.DueDate)
	if !ok {
		return
	}
	reqTs, ok := parseRequestTimestamp(w, req.RequestTimestamp)
	if !ok {
		return
	}
	created, err := h.taskUC.Create(r.Context(), user.UserID, task.TaskInput{
		Title:            req.Title,
		Content:          req.Content,
		DueDate:          dueDate,
		RequestTimestamp: reqTs,
	})
	if err != nil {
		appHttp.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	appHttp.JSON(w, http.StatusCreated, toTaskResponse(created))
}

// List handles GET /tasks.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := userOr401(w, r)
	if !ok {
		return
	}
	tasks, err := h.taskUC.List(r.Context(), user.UserID)
	if err != nil {
		appHttp.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resp := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		resp[i] = toTaskResponse(t)
	}
	appHttp.JSON(w, http.StatusOK, resp)
}

// Get handles GET /tasks/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	user, ok := userOr401(w, r)
	if !ok {
		return
	}
	id := chi.URLParam(r, "id")
	task, err := h.taskUC.Get(r.Context(), user.UserID, id)
	if err != nil {
		if appHttp.WritePortError(w, err) {
			return
		}
		appHttp.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	appHttp.JSON(w, http.StatusOK, toTaskResponse(task))
}

// Update handles PUT /tasks/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user, ok := userOr401(w, r)
	if !ok {
		return
	}
	id := chi.URLParam(r, "id")
	var req UpdateTaskRequest
	if !appHttp.ReadAndValidateJSON(w, r, &req) {
		return
	}
	reqTs, ok := parseRequestTimestamp(w, req.RequestTimestamp)
	if !ok {
		return
	}
	in := task.TaskUpdateInput{RequestTimestamp: reqTs}
	if req.Title != nil {
		in.Title = req.Title
	}
	if req.Content != nil {
		in.Content = req.Content
	}
	if req.DueDate != nil {
		d, ok := parseDueDate(w, *req.DueDate)
		if !ok {
			return
		}
		in.DueDate = &d
	}
	if req.Done != nil {
		in.Done = req.Done
	}
	t, err := h.taskUC.Update(r.Context(), user.UserID, id, in)
	if err != nil {
		if appHttp.WritePortError(w, err) {
			return
		}
		appHttp.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	appHttp.JSON(w, http.StatusOK, toTaskResponse(t))
}

// Delete handles DELETE /tasks/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := userOr401(w, r)
	if !ok {
		return
	}
	id := chi.URLParam(r, "id")
	var req DeleteTaskRequest
	if !appHttp.ReadAndValidateJSON(w, r, &req) {
		return
	}
	reqTs, ok := parseRequestTimestamp(w, req.RequestTimestamp)
	if !ok {
		return
	}
	err := h.taskUC.Delete(r.Context(), user.UserID, id, reqTs)
	if err != nil {
		if appHttp.WritePortError(w, err) {
			return
		}
		appHttp.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
}
