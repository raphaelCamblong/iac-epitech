package task

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/task-manager/api/internal/entity"
	"github.com/task-manager/api/internal/port"
)

type mockTaskRepo struct {
	create    func(ctx context.Context, task *entity.Task) error
	getByID   func(ctx context.Context, id, userID string) (*entity.Task, error)
	list      func(ctx context.Context, userID string) ([]*entity.Task, error)
	update    func(ctx context.Context, task *entity.Task) error
	delete    func(ctx context.Context, id, userID string) error
}

func (m *mockTaskRepo) Create(ctx context.Context, task *entity.Task) error {
	if m.create != nil {
		return m.create(ctx, task)
	}
	return nil
}
func (m *mockTaskRepo) GetByID(ctx context.Context, id, userID string) (*entity.Task, error) {
	if m.getByID != nil {
		return m.getByID(ctx, id, userID)
	}
	return nil, nil
}
func (m *mockTaskRepo) List(ctx context.Context, userID string) ([]*entity.Task, error) {
	if m.list != nil {
		return m.list(ctx, userID)
	}
	return nil, nil
}
func (m *mockTaskRepo) Update(ctx context.Context, task *entity.Task) error {
	if m.update != nil {
		return m.update(ctx, task)
	}
	return nil
}
func (m *mockTaskRepo) Delete(ctx context.Context, id, userID string) error {
	if m.delete != nil {
		return m.delete(ctx, id, userID)
	}
	return nil
}

func TestUseCase_Create(t *testing.T) {
	uc := NewTaskUseCase(&mockTaskRepo{
		create: func(ctx context.Context, task *entity.Task) error { return nil },
	})
	task, err := uc.Create(context.Background(), "user1", TaskInput{
		Title:            "Test",
		Content:          "Content",
		DueDate:          time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
		RequestTimestamp: time.Now(),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, task.ID)
	assert.Equal(t, "user1", task.UserID)
	assert.Equal(t, "Test", task.Title)
	assert.Equal(t, "Content", task.Content)
	assert.False(t, task.Done)
}

func TestUseCase_Get_NotFound(t *testing.T) {
	uc := NewTaskUseCase(&mockTaskRepo{
		getByID: func(ctx context.Context, id, userID string) (*entity.Task, error) { return nil, nil },
	})
	task, err := uc.Get(context.Background(), "user1", "nonexistent")
	assert.ErrorIs(t, err, port.ErrNotFound)
	assert.Nil(t, task)
}

func TestUseCase_Get_Found(t *testing.T) {
	expected := &entity.Task{ID: "t1", UserID: "u1", Title: "T"}
	uc := NewTaskUseCase(&mockTaskRepo{
		getByID: func(ctx context.Context, id, userID string) (*entity.Task, error) { return expected, nil },
	})
	task, err := uc.Get(context.Background(), "u1", "t1")
	require.NoError(t, err)
	assert.Equal(t, expected, task)
}

func TestUseCase_Update_Conflict(t *testing.T) {
	existing := &entity.Task{
		ID:               "t1",
		UserID:           "u1",
		RequestTimestamp: time.Date(2025, 9, 26, 10, 0, 0, 0, time.UTC),
	}
	uc := NewTaskUseCase(&mockTaskRepo{
		getByID: func(ctx context.Context, id, userID string) (*entity.Task, error) { return existing, nil },
	})
	task, err := uc.Update(context.Background(), "u1", "t1", TaskUpdateInput{
		Title:            strPtr("New"),
		RequestTimestamp: time.Date(2025, 9, 25, 10, 0, 0, 0, time.UTC),
	})
	assert.ErrorIs(t, err, port.ErrConflict)
	assert.Nil(t, task)
}

func TestUseCase_Update_Success(t *testing.T) {
	existing := &entity.Task{
		ID:               "t1",
		UserID:           "u1",
		Title:            "Old",
		RequestTimestamp: time.Date(2025, 9, 25, 10, 0, 0, 0, time.UTC),
	}
	var updated *entity.Task
	uc := NewTaskUseCase(&mockTaskRepo{
		getByID: func(ctx context.Context, id, userID string) (*entity.Task, error) { return existing, nil },
		update:  func(ctx context.Context, task *entity.Task) error { updated = task; return nil },
	})
	task, err := uc.Update(context.Background(), "u1", "t1", TaskUpdateInput{
		Title:            strPtr("New"),
		RequestTimestamp: time.Date(2025, 9, 26, 10, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	assert.Equal(t, "New", task.Title)
	assert.Equal(t, "New", updated.Title)
}

func TestUseCase_Delete_Conflict(t *testing.T) {
	existing := &entity.Task{
		ID:               "t1",
		RequestTimestamp: time.Date(2025, 9, 26, 10, 0, 0, 0, time.UTC),
	}
	uc := NewTaskUseCase(&mockTaskRepo{
		getByID: func(ctx context.Context, id, userID string) (*entity.Task, error) { return existing, nil },
	})
	err := uc.Delete(context.Background(), "u1", "t1", time.Date(2025, 9, 25, 10, 0, 0, 0, time.UTC))
	assert.ErrorIs(t, err, port.ErrConflict)
}

func TestUseCase_Delete_NotFound(t *testing.T) {
	uc := NewTaskUseCase(&mockTaskRepo{
		getByID: func(ctx context.Context, id, userID string) (*entity.Task, error) { return nil, nil },
	})
	err := uc.Delete(context.Background(), "u1", "t1", time.Now())
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestUseCase_Delete_Success(t *testing.T) {
	existing := &entity.Task{ID: "t1", RequestTimestamp: time.Date(2025, 9, 25, 10, 0, 0, 0, time.UTC)}
	deleted := false
	uc := NewTaskUseCase(&mockTaskRepo{
		getByID: func(ctx context.Context, id, userID string) (*entity.Task, error) { return existing, nil },
		delete:  func(ctx context.Context, id, userID string) error { deleted = true; return nil },
	})
	err := uc.Delete(context.Background(), "u1", "t1", time.Date(2025, 9, 26, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestUseCase_Get_RepoError(t *testing.T) {
	uc := NewTaskUseCase(&mockTaskRepo{
		getByID: func(ctx context.Context, id, userID string) (*entity.Task, error) {
			return nil, errors.New("db error")
		},
	})
	_, err := uc.Get(context.Background(), "u1", "t1")
	assert.Error(t, err)
}

func strPtr(s string) *string { return &s }
