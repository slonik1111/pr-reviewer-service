package service

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	repo "github.com/slonik1111/pr-reviewer-service/internal/repository"
)

var (
	ErrTeamNotFound = errors.New("team not found")
)

type TeamService struct {
	userRepo repo.UserRepository
}

func NewTeamService(userRepo repo.UserRepository) *TeamService {
	return &TeamService{
		userRepo: userRepo,
	}
}

func (s *TeamService) CreateTeam(team string, members []domain.User) error {
	_, err := s.userRepo.ListTeamUsers(team)
	if err == nil {
		return fmt.Errorf("team %s already exists", team)
	}
	for _, u := range members {
		u.TeamName = team
		if err := s.userRepo.Create(u); err != nil {
			_ = s.userRepo.Update(u)
		}
	}

	return nil
}

func (s *TeamService) GetTeam(teamName string) ([]domain.User, error) {
	team, err := s.userRepo.ListTeamUsers(teamName)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	return team, nil
}

func (s *TeamService) ListActiveMembers(teamID string) ([]domain.User, error) {
	users, err := s.userRepo.ListActiveTeamUsers(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}
	return users, nil
}

func (s *TeamService) GetRandomActiveMembers(teamID string, excludeIDs []string, n int) ([]domain.User, error) {
	users, err := s.ListActiveMembers(teamID)
	if err != nil {
		return nil, err
	}

	excludeMap := make(map[string]bool)
	for _, id := range excludeIDs {
		excludeMap[id] = true
	}

	filtered := []domain.User{}
	for _, u := range users {
		if !excludeMap[u.ID] {
			filtered = append(filtered, u)
		}
	}

	rand.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})

	if len(filtered) > n {
		filtered = filtered[:n]
	}

	return filtered, nil
}
