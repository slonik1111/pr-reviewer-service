package inmemory

import (
	"errors"
	"sync"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	"github.com/slonik1111/pr-reviewer-service/internal/repository"
)

var _ repo.UserRepository = (*UserRepoInMemory)(nil)

type UserRepoInMemory struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewUserRepoInMemory() *UserRepoInMemory {
	return &UserRepoInMemory{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepoInMemory) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	r.users[user.ID] = user
	return nil
}

func (r *UserRepoInMemory) GetUser(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *UserRepoInMemory) Update(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	r.users[user.ID] = user
	return nil
}

func (r *UserRepoInMemory) ListTeamUsers(teamID string) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := []*domain.User{}
	for _, u := range r.users {
		if u.TeamID == teamID {
			res = append(res, u)
		}
	}
	return res, nil
}

func (r *UserRepoInMemory) ListActiveTeamUsers(teamID string) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := []*domain.User{}
	for _, u := range r.users {
		if u.TeamID == teamID && u.IsActive {
			res = append(res, u)
		}
	}
	return res, nil
}
