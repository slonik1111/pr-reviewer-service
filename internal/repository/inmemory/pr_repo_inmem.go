package inmemory

import (
	"errors"
	"sync"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	"github.com/slonik1111/pr-reviewer-service/internal/repository"
)

var _ repo.PullRequestRepository = nil

type PRRepoInMemory struct {
	mu  sync.RWMutex
	prs map[string]domain.PullRequest
}

func NewPRRepoInMemory() *PRRepoInMemory {
	return &PRRepoInMemory{
		prs: make(map[string]domain.PullRequest),
	}
}

func (r *PRRepoInMemory) Create(pr domain.PullRequest) (domain.PullRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prs[pr.ID]; exists {
		return domain.PullRequest{}, errors.New("PR already exists")
	}
	r.prs[pr.ID] = pr
	return pr, nil
}

func (r *PRRepoInMemory) Get(id string) (domain.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pr, ok := r.prs[id]
	if !ok {
		return domain.PullRequest{}, errors.New("PR not found")
	}
	return pr, nil
}

func (r *PRRepoInMemory) Update(pr domain.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prs[pr.ID]; !exists {
		return errors.New("PR not found")
	}
	r.prs[pr.ID] = pr
	return nil
}

func (r *PRRepoInMemory) ListByReviewer(userID string) ([]domain.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := []domain.PullRequest{}
	for _, pr := range r.prs {
		if (len(pr.Reviewers) != 0 && pr.Reviewers[0] == userID) || (len(pr.Reviewers) > 1 && pr.Reviewers[1] == userID) {
			res = append(res, pr)
		}
	}
	return res, nil
}

func (r *PRRepoInMemory) ListByAuthor(userID string) ([]domain.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := []domain.PullRequest{}
	for _, pr := range r.prs {
		if pr.AuthorID == userID {
			res = append(res, pr)
		}
	}
	return res, nil
}
