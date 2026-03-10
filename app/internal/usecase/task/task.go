package task

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/task-manager/api/internal/entity"
	"github.com/task-manager/api/internal/port"
)

// TaskInput represents input for creating a task.
type TaskInput struct {
	Title            string
	Content          string
	DueDate          time.Time
	RequestTimestamp time.Time
}

// TaskUpdateInput represents input for updating a task.
type TaskUpdateInput struct {
	Title            *string
	Content          *string
	DueDate          *time.Time
	Done             *bool
	RequestTimestamp time.Time
}

// UseCase handles task business logic.
type UseCase struct {
	repo port.TaskRepository
}

// NewTaskUseCase creates a new task UseCase.
func NewTaskUseCase(repo port.TaskRepository) *UseCase {
	return &UseCase{repo: repo}
}

// Create creates a new task for the given user.
func (uc *UseCase) Create(ctx context.Context, userID string, in TaskInput) (*entity.Task, error) {
	now := time.Now()
	task := &entity.Task{
		ID:               uuid.New().String(),
		UserID:           userID,
		Title:            in.Title,
		Content:          in.Content,
		DueDate:          in.DueDate,
		Done:             false,
		RequestTimestamp: in.RequestTimestamp,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := uc.repo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}

// Get retrieves a task by ID for the given user.
func (uc *UseCase) Get(ctx context.Context, userID, id string) (*entity.Task, error) {
	task, err := uc.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	if task == nil {
		return nil, port.ErrNotFound
	}
	return task, nil
}

// List lists all tasks for the given user.
func (uc *UseCase) List(ctx context.Context, userID string) ([]*entity.Task, error) {
	tasks, err := uc.repo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return tasks, nil
}

// Update updates a task. Returns ErrConflict if request_timestamp is <= stored last timestamp.
func (uc *UseCase) Update(ctx context.Context, userID, id string, in TaskUpdateInput) (*entity.Task, error) {
	existing, err := uc.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("get task for update: %w", err)
	}
	if existing == nil {
		return nil, port.ErrNotFound
	}

	if in.RequestTimestamp.Before(existing.RequestTimestamp) || in.RequestTimestamp.Equal(existing.RequestTimestamp) {
		return nil, port.ErrConflict
	}

	updated := &entity.Task{
		ID:               existing.ID,
		UserID:           existing.UserID,
		Title:            existing.Title,
		Content:          existing.Content,
		DueDate:          existing.DueDate,
		Done:             existing.Done,
		RequestTimestamp: in.RequestTimestamp,
		CreatedAt:        existing.CreatedAt,
		UpdatedAt:        time.Now(),
	}
	if in.Title != nil {
		updated.Title = *in.Title
	}
	if in.Content != nil {
		updated.Content = *in.Content
	}
	if in.DueDate != nil {
		updated.DueDate = *in.DueDate
	}
	if in.Done != nil {
		updated.Done = *in.Done
	}

	if err := uc.repo.Update(ctx, updated); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}
	return updated, nil
}

// Delete deletes a task. Returns ErrConflict if request_timestamp is <= stored last timestamp.
func (uc *UseCase) Delete(ctx context.Context, userID, id string, requestTimestamp time.Time) error {
	existing, err := uc.repo.GetByID(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("get task for delete: %w", err)
	}
	if existing == nil {
		return port.ErrNotFound
	}

	if requestTimestamp.Before(existing.RequestTimestamp) || requestTimestamp.Equal(existing.RequestTimestamp) {
		return port.ErrConflict
	}

	if err := uc.repo.Delete(ctx, id, userID); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}
