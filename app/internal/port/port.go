package port

import (
	"context"

	"github.com/task-manager/api/internal/entity"
)

// TaskRepository defines the interface for task persistence.
type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id, userID string) (*entity.Task, error)
	List(ctx context.Context, userID string) ([]*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id, userID string) error
}

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
