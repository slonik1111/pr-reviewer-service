package postgres

import (
	"database/sql"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	"github.com/slonik1111/pr-reviewer-service/internal/repository"
	"github.com/lib/pq"
)

var _ repo.PullRequestRepository = (*PullRequestRepo)(nil)

type PullRequestRepo struct {
	db *sql.DB
}

func NewPullRequestRepo(db *sql.DB) *PullRequestRepo {
	return &PullRequestRepo{db: db}
}

func (r *PullRequestRepo) Create(pr domain.PullRequest) (domain.PullRequest, error) {
	_, err := r.db.Exec(`
		INSERT INTO pull_requests 
			(id, title, description, author_id, team_id, reviewers, status, merged_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, pr.ID, pr.Title, pr.Description, pr.AuthorID, pr.TeamID,
		pq.Array(pr.Reviewers), pr.Status, pr.MergedAt)

	return pr, err
}

func (r *PullRequestRepo) Get(id string) (domain.PullRequest, error) {
	var pr domain.PullRequest
	var reviewers []string
	var mergedAt sql.NullString

	err := r.db.QueryRow(`
		SELECT id, title, description, author_id, team_id, reviewers, status, merged_at
		FROM pull_requests
		WHERE id = $1
	`, id).Scan(
		&pr.ID, &pr.Title, &pr.Description, &pr.AuthorID, &pr.TeamID,
		pq.Array(&reviewers), &pr.Status, &mergedAt,
	)

	if err != nil {
		return domain.PullRequest{}, err
	}

	pr.Reviewers = reviewers

	if mergedAt.Valid {
		v := mergedAt.String
		pr.MergedAt = &v
	}

	return pr, nil
}


func (r *PullRequestRepo) Update(pr domain.PullRequest) error {
	_, err := r.db.Exec(`
		UPDATE pull_requests
		SET title = $1, description = $2,
		    author_id = $3, team_id = $4,
		    reviewers = $5, status = $6, merged_at = $7
		WHERE id = $8
	`, pr.Title, pr.Description, pr.AuthorID, pr.TeamID,
		pq.Array(pr.Reviewers), pr.Status, pr.MergedAt, pr.ID)

	return err
}

func (r *PullRequestRepo) ListByReviewer(userID string) ([]domain.PullRequest, error) {
	rows, err := r.db.Query(`
		SELECT id, title, description, author_id, team_id,
		       reviewers, status, merged_at
		FROM pull_requests
		WHERE $1 = ANY(reviewers)
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.PullRequest

	for rows.Next() {
		var pr domain.PullRequest
		var reviewers []string
		var mergedAt sql.NullString

		if err := rows.Scan(
			&pr.ID, &pr.Title, &pr.Description,
			&pr.AuthorID, &pr.TeamID,
			pq.Array(&reviewers),
			&pr.Status, &mergedAt,
		); err != nil {
			return nil, err
		}

  pr.Reviewers = reviewers

		if mergedAt.Valid {
			v := mergedAt.String
			pr.MergedAt = &v
		}

		list = append(list, pr)
	}

	return list, nil
}

func (r *PullRequestRepo) ListByAuthor(userID string) ([]domain.PullRequest, error) {
	rows, err := r.db.Query(`
		SELECT id, title, description, author_id, team_id,
		       reviewers, status, merged_at
		FROM pull_requests
		WHERE author_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.PullRequest

	for rows.Next() {
		var pr domain.PullRequest
		var reviewers []string
		var mergedAt sql.NullString

		if err := rows.Scan(
			&pr.ID, &pr.Title, &pr.Description,
			&pr.AuthorID, &pr.TeamID,
			pq.Array(&reviewers),
			&pr.Status, &mergedAt,
		); err != nil {
			return nil, err
		}

		pr.Reviewers = reviewers

		if mergedAt.Valid {
			v := mergedAt.String
			pr.MergedAt = &v
		}

		list = append(list, pr)
	}

	return list, nil
}
