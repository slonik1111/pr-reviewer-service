package postgres

import (
	"database/sql"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	"github.com/slonik1111/pr-reviewer-service/internal/repository"
)

var _ repo.UserRepository = (*UserRepo)(nil)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user domain.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (id, username, team_id, is_active)
		VALUES ($1, $2, $3, $4)
	`, user.ID, user.Name, user.TeamName, user.IsActive)
	return err
}

func (r *UserRepo) GetUser(id string) (domain.User, error) {
	var u domain.User

	err := r.db.QueryRow(`
		SELECT id, username, team_id, is_active
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive)

	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserRepo) Update(user domain.User) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET username = $1, team_id = $2, is_active = $3
		WHERE id = $4
	`, user.Name, user.TeamName, user.IsActive, user.ID)
	return err
}

func (r *UserRepo) ListTeamUsers(teamID string) ([]domain.User, error) {
	rows, err := r.db.Query(`
		SELECT id, username, team_id, is_active
		FROM users WHERE team_id = $1
	`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.User

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		list = append(list, u)
	}

	return list, nil
}

func (r *UserRepo) ListActiveTeamUsers(teamID string) ([]domain.User, error) {
	rows, err := r.db.Query(`
		SELECT id, username, team_id, is_active
		FROM users 
		WHERE team_id = $1 AND is_active = true
	`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.User

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		list = append(list, u)
	}

	return list, nil
}
