package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/task-manager/api/internal/entity"
	"github.com/task-manager/api/internal/port"
)

// AuthInput represents input for registration.
type AuthInput struct {
	Email    string
	Password string
}

// UseCase handles authentication business logic.
type UseCase struct {
	userRepo port.UserRepository
}

// NewAuthUseCase creates a new auth UseCase.
func NewAuthUseCase(userRepo port.UserRepository) *UseCase {
	return &UseCase{userRepo: userRepo}
}

// Register creates a new user and returns the user entity.
func (uc *UseCase) Register(ctx context.Context, in AuthInput) (*entity.User, error) {
	existing, err := uc.userRepo.GetByEmail(ctx, in.Email)
	if err != nil {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return nil, port.ErrDuplicateUser
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &entity.User{
		ID:           uuid.New().String(),
		Email:        in.Email,
		PasswordHash: string(hash),
	}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

// Login validates credentials and returns the user if valid.
func (uc *UseCase) Login(ctx context.Context, in AuthInput) (*entity.User, error) {
	user, err := uc.userRepo.GetByEmail(ctx, in.Email)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, port.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, port.ErrUnauthorized
	}
	return user, nil
}
