package repo

import "github.com/slonik1111/pr-reviewer-service/internal/domain"

type PullRequestRepository interface {
	Create(pr domain.PullRequest) (domain.PullRequest, error)

	Get(id string) (domain.PullRequest, error)

	Update(pr domain.PullRequest) error

	ListByReviewer(userID string) ([]domain.PullRequest, error)

	ListByAuthor(userID string) ([]domain.PullRequest, error)
}
