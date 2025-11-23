package service

import (
	"errors"
	"fmt"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	repo "github.com/slonik1111/pr-reviewer-service/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	userRepo repo.UserRepository
	prRepo   repo.PullRequestRepository
}

func NewUserService(userRepo repo.UserRepository, prRepo repo.PullRequestRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (s *UserService) GetUser(userID string) (domain.User, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return domain.User{}, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) SetIsActive(userID string, isActive bool) (domain.User, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return domain.User{}, ErrUserNotFound
	}

	user.IsActive = isActive
	if err := s.userRepo.Update(user); err != nil {
		return domain.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *UserService) UpsertUser(user domain.User) error {
	existing, err := s.userRepo.GetUser(user.ID)
	if err != nil {
		return s.userRepo.Create(user)
	}

	existing.Name = user.Name
	existing.TeamName = user.TeamName
	existing.IsActive = user.IsActive

	return s.userRepo.Update(existing)
}

func (s *UserService) GetPRsForReview(userID string) ([]domain.PullRequest, error) {
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
