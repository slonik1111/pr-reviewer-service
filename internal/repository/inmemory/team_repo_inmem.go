package inmemory

import (
	"errors"
	"sync"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	"github.com/slonik1111/pr-reviewer-service/internal/repository"
)

var _ repo.TeamRepository = (*TeamRepoInMemory)(nil)

type TeamRepoInMemory struct {
	mu    sync.RWMutex
	teams map[string]*domain.Team
}

func NewTeamRepoInMemory() *TeamRepoInMemory {
	return &TeamRepoInMemory{
		teams: make(map[string]*domain.Team),
	}
}

func (r *TeamRepoInMemory) Create(team *domain.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.teams[team.Name]; exists {
		return errors.New("team already exists")
	}
	r.teams[team.Name] = team
	return nil
}

func (r *TeamRepoInMemory) GetTeam(id string) (*domain.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.teams[id]
	if !ok {
		return nil, errors.New("team not found")
	}
	return t, nil
}

func (r *TeamRepoInMemory) GetTeamByName(name string) (*domain.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, t := range r.teams {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, errors.New("team not found")
}

func (r *TeamRepoInMemory) AddUser(teamID string, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.teams[teamID]
	if !ok {
		return errors.New("team not found")
	}
	t.Users = append(t.Users, userID)
	return nil
}

func (r *TeamRepoInMemory) ListTeams() ([]*domain.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := []*domain.Team{}
	for _, t := range r.teams {
		res = append(res, t)
	}
	return res, nil
}
