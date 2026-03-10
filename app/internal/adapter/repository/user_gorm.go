package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/task-manager/api/internal/adapter/repository/model"
	"github.com/task-manager/api/internal/entity"
)

// UserGormRepository implements usecase.UserRepository using GORM.
type UserGormRepository struct {
	db *gorm.DB
}

// NewUserGormRepository creates a new UserGormRepository.
func NewUserGormRepository(db *gorm.DB) *UserGormRepository {
	return &UserGormRepository{db: db}
}

func toUserModel(u *entity.User) *model.UserModel {
	return &model.UserModel{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}
}

func toUserEntity(m *model.UserModel) *entity.User {
	return &entity.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
	}
}

// Create inserts a user.
func (r *UserGormRepository) Create(ctx context.Context, user *entity.User) error {
	m := toUserModel(user)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// GetByEmail retrieves a user by email.
func (r *UserGormRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.UserModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return toUserEntity(&m), nil
}
