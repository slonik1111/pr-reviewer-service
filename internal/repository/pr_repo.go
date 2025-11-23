package repo

import "github.com/slonik1111/pr-reviewer-service/internal/domain"

type PullRequestRepository interface {
	// Create a PR
	Create(pr domain.PullRequest) (domain.PullRequest, error)

	// Get PR by ID
	Get(id string) (domain.PullRequest, error)

	// Update PR (e.g. merge, reassign)
	Update(pr domain.PullRequest) error

	// List PRs assigned to reviewer
	ListByReviewer(userID string) ([]domain.PullRequest, error)

	// List PRs authored by user
	ListByAuthor(userID string) ([]domain.PullRequest, error)
}
