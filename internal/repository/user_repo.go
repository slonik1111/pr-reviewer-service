package repo

import "github.com/slonik1111/pr-reviewer-service/internal/domain"


type UserRepository interface {
	Create(user domain.User) error

	GetUser(id string) (domain.User, error)

	Update(user domain.User) error

	ListTeamUsers(teamID string) ([]domain.User, error)

	ListActiveTeamUsers(teamID string) ([]domain.User, error)
}
