package usecase

import (
	"backend/internal/config"
	"backend/internal/domain"
	"backend/internal/dto"
	"backend/internal/infrastructure"
	"backend/internal/repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthUsecase(userRepo repository.UserRepository, cfg *config.Config) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (u *AuthUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {

	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := infrastructure.GenerateJWT(
		user.ID,
		user.Email,
		u.cfg.JWTSecret,
		u.cfg.JWTExpireHour,
	)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
	}, nil
}

func (u *AuthUsecase) Signup(ctx context.Context, req dto.SignupRequest) (*dto.LoginResponse, error) {
	_, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := infrastructure.GenerateJWT(
		user.ID,
		user.Email,
		u.cfg.JWTSecret,
		u.cfg.JWTExpireHour,
	)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
	}, nil
}
