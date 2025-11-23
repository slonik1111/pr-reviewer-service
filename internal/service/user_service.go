package service

import (
	"errors"
	"fmt"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	"github.com/slonik1111/pr-reviewer-service/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// UserService содержит бизнес-логику для пользователей
type UserService struct {
	userRepo repo.UserRepository
	prRepo   repo.PullRequestRepository
}

// NewUserService создает новый сервис пользователей
func NewUserService(userRepo repo.UserRepository, prRepo repo.PullRequestRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

// GetUser возвращает пользователя по ID
func (s *UserService) GetUser(userID string) (*domain.User, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// SetIsActive устанавливает флаг активности пользователя
func (s *UserService) SetIsActive(userID string, isActive bool) (*domain.User, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.IsActive = isActive
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// UpsertUser создает нового пользователя или обновляет существующего
func (s *UserService) UpsertUser(user *domain.User) error {
	existing, err := s.userRepo.GetUser(user.ID)
	if err != nil {
		// если не найден, создаем
		return s.userRepo.Create(user)
	}

	// если найден — обновляем
	existing.Name = user.Name
	existing.TeamID = user.TeamID
	existing.IsActive = user.IsActive

	return s.userRepo.Update(existing)
}

// GetPRsForReview возвращает список PR, где пользователь назначен ревьювером
func (s *UserService) GetPRsForReview(userID string) ([]*domain.PullRequest, error) {
	// проверяем существование пользователя
	_, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	prs, err := s.prRepo.ListByReviewer(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list PRs for reviewer: %w", err)
	}

	return prs, nil
}
