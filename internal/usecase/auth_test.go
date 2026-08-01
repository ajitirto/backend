package usecase

import (
	"backend/internal/config"
	"backend/internal/domain"
	"backend/internal/dto"
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	findByEmail func(ctx context.Context, email string) (*domain.User, error)
	create      func(ctx context.Context, user *domain.User) error
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.findByEmail(ctx, email)
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.create(ctx, user)
}

func newConfig() *config.Config {
	return &config.Config{
		JWTSecret:     "secret",
		JWTExpireHour: 24,
	}
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	return string(hash)
}

func TestLogin_Success(t *testing.T) {
	repo := &mockUserRepository{
		findByEmail: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{
				ID:       1,
				Email:    email,
				Password: hashPassword("password123"),
			}, nil
		},
		create: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}

	uc := NewAuthUsecase(repo, newConfig())

	ctx := context.Background()
	resp, err := uc.Login(ctx, dto.LoginRequest{
		Email:    "admin@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response")
	}

	if resp.Token == "" {
		t.Fatal("expected token")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		findByEmail: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
		create: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}

	uc := NewAuthUsecase(repo, newConfig())

	ctx := context.Background()
	_, err := uc.Login(ctx, dto.LoginRequest{
		Email:    "admin@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := &mockUserRepository{
		findByEmail: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{
				ID:       1,
				Email:    email,
				Password: hashPassword("password123"),
			}, nil
		},
		create: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}

	uc := NewAuthUsecase(repo, newConfig())

	ctx := context.Background()
	_, err := uc.Login(ctx, dto.LoginRequest{
		Email:    "admin@example.com",
		Password: "wrong-password",
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSignup_Success(t *testing.T) {
	repo := &mockUserRepository{
		findByEmail: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
		create: func(ctx context.Context, user *domain.User) error {
			user.ID = 1
			return nil
		},
	}

	uc := NewAuthUsecase(repo, newConfig())

	ctx := context.Background()
	resp, err := uc.Signup(ctx, dto.SignupRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Token == "" {
		t.Fatal("expected token")
	}
}

func TestSignup_EmailAlreadyExists(t *testing.T) {
	repo := &mockUserRepository{
		findByEmail: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{
				ID:    1,
				Email: email,
			}, nil
		},
		create: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}

	uc := NewAuthUsecase(repo, newConfig())
	ctx := context.Background()
	_, err := uc.Signup(ctx, dto.SignupRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSignup_CreateFailed(t *testing.T) {
	repo := &mockUserRepository{
		findByEmail: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
		create: func(ctx context.Context, user *domain.User) error {
			return errors.New("database error")
		},
	}

	uc := NewAuthUsecase(repo, newConfig())

	ctx := context.Background()
	_, err := uc.Signup(ctx, dto.SignupRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error")
	}
}
