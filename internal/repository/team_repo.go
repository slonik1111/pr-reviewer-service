package repo

import "github.com/slonik1111/pr-reviewer-service/internal/domain"

type TeamRepository interface {
	// Create new team
	Create(team *domain.Team) error

	// Get team by ID
	GetTeam(id string) (*domain.Team, error)

	// Get team by name (team names must be unique)
	GetTeamByName(name string) (*domain.Team, error)

	// AddUser attaches a user ID to a team
	AddUser(teamID string, userID string) error

	// List all teams
	ListTeams() ([]*domain.Team, error)
}
