package repo

import "github.com/slonik1111/pr-reviewer-service/internal/domain"


type UserRepository interface {
	// Create user or return error if ID already exists
	Create(user *domain.User) error

	// GetUser returns user by id. Returns (nil, ErrNotFound) if user does not exist
	GetUser(id string) (*domain.User, error)

	// Update modifies fields of existing user
	Update(user *domain.User) error

	// ListTeamUsers returns all users that belong to a team
	ListTeamUsers(teamID string) ([]*domain.User, error)

	// ListActiveTeamUsers returns only active users of a team
	ListActiveTeamUsers(teamID string) ([]*domain.User, error)
}
