package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/task-manager/api/internal/adapter/repository/model"
	"github.com/task-manager/api/internal/entity"
	"github.com/task-manager/api/internal/port"
)

// TaskGormRepository implements usecase.TaskRepository using GORM.
type TaskGormRepository struct {
	db *gorm.DB
}

// NewTaskGormRepository creates a new TaskGormRepository.
func NewTaskGormRepository(db *gorm.DB) *TaskGormRepository {
	return &TaskGormRepository{db: db}
}

func toTaskModel(t *entity.Task) *model.TaskModel {
	return &model.TaskModel{
		ID:               t.ID,
		UserID:           t.UserID,
		Title:            t.Title,
		Content:          t.Content,
		DueDate:          t.DueDate,
		Done:             t.Done,
		RequestTimestamp: t.RequestTimestamp,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

func toTaskEntity(m *model.TaskModel) *entity.Task {
	return &entity.Task{
		ID:               m.ID,
		UserID:           m.UserID,
		Title:            m.Title,
		Content:          m.Content,
		DueDate:          m.DueDate,
		Done:             m.Done,
		RequestTimestamp: m.RequestTimestamp,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

// Create inserts a task.
func (r *TaskGormRepository) Create(ctx context.Context, task *entity.Task) error {
	m := toTaskModel(task)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

// GetByID retrieves a task by ID and user ID.
func (r *TaskGormRepository) GetByID(ctx context.Context, id, userID string) (*entity.Task, error) {
	var m model.TaskModel
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get task: %w", err)
	}
	return toTaskEntity(&m), nil
}

// List returns all tasks for a user.
func (r *TaskGormRepository) List(ctx context.Context, userID string) ([]*entity.Task, error) {
	var models []model.TaskModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	tasks := make([]*entity.Task, len(models))
	for i := range models {
		tasks[i] = toTaskEntity(&models[i])
	}
	return tasks, nil
}

// Update updates a task.
func (r *TaskGormRepository) Update(ctx context.Context, task *entity.Task) error {
	result := r.db.WithContext(ctx).Model(&model.TaskModel{}).
		Where("id = ? AND user_id = ?", task.ID, task.UserID).
		Updates(map[string]interface{}{
			"title":             task.Title,
			"content":           task.Content,
			"due_date":          task.DueDate,
			"done":              task.Done,
			"request_timestamp": task.RequestTimestamp,
			"updated_at":        task.UpdatedAt,
		})
	if result.Error != nil {
		return fmt.Errorf("update task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return port.ErrNotFound
	}
	return nil
}

// Delete removes a task.
func (r *TaskGormRepository) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.TaskModel{})
	if result.Error != nil {
		return fmt.Errorf("delete task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return port.ErrNotFound
	}
	return nil
}
