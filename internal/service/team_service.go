package service

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	"github.com/slonik1111/pr-reviewer-service/internal/repository"
)

var (
	ErrTeamNotFound = errors.New("team not found")
)

// TeamService управляет командами и их участниками
type TeamService struct {
	teamRepo repo.TeamRepository
	userRepo repo.UserRepository
}

// NewTeamService создает сервис команд
func NewTeamService(teamRepo repo.TeamRepository, userRepo repo.UserRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// CreateTeam создает команду и добавляет участников
func (s *TeamService) CreateTeam(team *domain.Team, members []*domain.User) error {
	// проверяем уникальность команды
	_, err := s.teamRepo.GetTeamByName(team.Name)
	if err == nil {
		return fmt.Errorf("team %s already exists", team.Name)
	}

	// создаем команду
	if err := s.teamRepo.Create(team); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	// добавляем пользователей и привязываем к команде
	for _, u := range members {
		u.TeamID = team.Name
		if err := s.userRepo.Create(u); err != nil {
			// если пользователь существует — обновим
			_ = s.userRepo.Update(u)
		}
		_ = s.teamRepo.AddUser(team.Name, u.ID)
	}

	return nil
}

// GetTeam возвращает команду по имени
func (s *TeamService) GetTeam(teamName string) (*domain.Team, error) {
	team, err := s.teamRepo.GetTeamByName(teamName)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	return team, nil
}

// ListActiveMembers возвращает список активных пользователей команды
func (s *TeamService) ListActiveMembers(teamID string) ([]*domain.User, error) {
	users, err := s.userRepo.ListActiveTeamUsers(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}
	return users, nil
}

// GetRandomActiveMembers возвращает до N случайных активных пользователей
func (s *TeamService) GetRandomActiveMembers(teamID string, excludeIDs []string, n int) ([]*domain.User, error) {
	users, err := s.ListActiveMembers(teamID)
	log.Println("ative users:", users)
	if err != nil {
		return nil, err
	}

	// фильтруем исключения
	excludeMap := make(map[string]bool)
	for _, id := range excludeIDs {
		excludeMap[id] = true
	}

	filtered := []*domain.User{}
	for _, u := range users {
		if !excludeMap[u.ID] {
			filtered = append(filtered, u)
		}
	}

	// перемешиваем
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})

	if len(filtered) > n {
		filtered = filtered[:n]
	}

	return filtered, nil
}
