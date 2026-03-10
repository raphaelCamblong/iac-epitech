package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/task-manager/api/internal/entity"
	"github.com/task-manager/api/internal/port"
)

type mockUserRepo struct {
	create     func(ctx context.Context, user *entity.User) error
	getByEmail func(ctx context.Context, email string) (*entity.User, error)
}

func (m *mockUserRepo) Create(ctx context.Context, user *entity.User) error {
	if m.create != nil {
		return m.create(ctx, user)
	}
	return nil
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.getByEmail != nil {
		return m.getByEmail(ctx, email)
	}
	return nil, nil
}

func TestUseCase_Register_Success(t *testing.T) {
	var created *entity.User
	uc := NewAuthUseCase(&mockUserRepo{
		getByEmail: func(ctx context.Context, email string) (*entity.User, error) { return nil, nil },
		create:     func(ctx context.Context, user *entity.User) error { created = user; return nil },
	})
	user, err := uc.Register(context.Background(), AuthInput{
		Email:    "test@example.com",
		Password: "password123",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotEmpty(t, user.PasswordHash)
	assert.Equal(t, user.ID, created.ID)
}

func TestUseCase_Register_Duplicate(t *testing.T) {
	uc := NewAuthUseCase(&mockUserRepo{
		getByEmail: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "u1", Email: email}, nil
		},
	})
	user, err := uc.Register(context.Background(), AuthInput{
		Email:    "test@example.com",
		Password: "password123",
	})
	assert.ErrorIs(t, err, port.ErrDuplicateUser)
	assert.Nil(t, user)
}

func TestUseCase_Login_Success(t *testing.T) {
	existing := &entity.User{
		ID:           "u1",
		Email:        "test@example.com",
		PasswordHash: mustHash("password123"),
	}
	uc := NewAuthUseCase(&mockUserRepo{
		getByEmail: func(ctx context.Context, email string) (*entity.User, error) { return existing, nil },
	})
	user, err := uc.Login(context.Background(), AuthInput{
		Email:    "test@example.com",
		Password: "password123",
	})
	require.NoError(t, err)
	assert.Equal(t, "u1", user.ID)
}

func TestUseCase_Login_InvalidPassword(t *testing.T) {
	existing := &entity.User{
		ID:           "u1",
		Email:        "test@example.com",
		PasswordHash: mustHash("password123"),
	}
	uc := NewAuthUseCase(&mockUserRepo{
		getByEmail: func(ctx context.Context, email string) (*entity.User, error) { return existing, nil },
	})
	user, err := uc.Login(context.Background(), AuthInput{
		Email:    "test@example.com",
		Password: "wrong",
	})
	assert.ErrorIs(t, err, port.ErrUnauthorized)
	assert.Nil(t, user)
}

func TestUseCase_Login_UserNotFound(t *testing.T) {
	uc := NewAuthUseCase(&mockUserRepo{
		getByEmail: func(ctx context.Context, email string) (*entity.User, error) { return nil, nil },
	})
	user, err := uc.Login(context.Background(), AuthInput{
		Email:    "test@example.com",
		Password: "password123",
	})
	assert.ErrorIs(t, err, port.ErrUnauthorized)
	assert.Nil(t, user)
}

func mustHash(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
	if err != nil {
		panic(err)
	}
	return string(hash)
}
