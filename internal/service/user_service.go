package service

import (
	"context"
	"errors"
	"strings"

	"github.com/sahilrana27582/payment-system/internal/repository"
)

var (
	ErrEmailTaken      = errors.New("email already in use")
	ErrPasswordTooLong = errors.New("password must be at most 72 bytes")
)

type UserRepository interface {
	CreateUser(ctx context.Context, username, password, email string) (string, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, username, password, email string) (string, error) {

	email = strings.ToLower(email)

	userID, err := s.repo.CreateUser(ctx, username, password, email)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return "", ErrEmailTaken
		}
		return "", err
	}

	return userID, nil
}
