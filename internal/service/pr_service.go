package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	repo "github.com/slonik1111/pr-reviewer-service/internal/repository"
)

var (
	ErrPRNotFound          = errors.New("pull request not found")
	ErrPRMerged            = errors.New("cannot modify a merged PR")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
)

// PRService управляет Pull Request и назначением ревьюверов
type PRService struct {
	prRepo   repo.PullRequestRepository
	userRepo repo.UserRepository
}

// NewPRService создает сервис PR
func NewPRService(prRepo repo.PullRequestRepository, userRepo repo.UserRepository) *PRService {
	return &PRService{
		prRepo:   prRepo,
		userRepo: userRepo,
	}
}

// CreatePR создает новый Pull Request и назначает до 2 ревьюверов
func (s *PRService) CreatePR(pr domain.PullRequest) (domain.PullRequest, error) {
	_, err := s.prRepo.Get(pr.ID)
	if err == nil {
		return domain.PullRequest{}, fmt.Errorf("PR %s already exists", pr.ID)
	}

	author, err := s.userRepo.GetUser(pr.AuthorID)
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("author %s not found", pr.AuthorID)
	}

	pr.TeamID = author.TeamName

	teamService := NewTeamService(s.userRepo)
	reviewers, _ := teamService.GetRandomActiveMembers(pr.TeamID, []string{pr.AuthorID}, 2)

	pr.Reviewers = make([]string, 0)

	if len(reviewers) > 0 {
		pr.Reviewers = append(pr.Reviewers, reviewers[0].ID)
	}
	if len(reviewers) > 1 {
		pr.Reviewers = append(pr.Reviewers, reviewers[1].ID)
	}

	pr.Status = domain.PRStatusOpen

	return s.prRepo.Create(pr)
}

func (s *PRService) MergePR(prID string) (domain.PullRequest, error) {
	pr, err := s.prRepo.Get(prID)
	if err != nil {
		return domain.PullRequest{}, ErrPRNotFound
	}

	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	pr.Status = domain.PRStatusMerged
	merged := time.Now().UTC().Format(time.RFC3339)
	pr.MergedAt = &merged

	if err := s.prRepo.Update(pr); err != nil {
		return domain.PullRequest{}, fmt.Errorf("failed to merge PR: %w", err)
	}

	return pr, nil
}

// ReassignReviewer заменяет одного ревьювера на другого из команды
func (s *PRService) ReassignReviewer(prID string, oldReviewerID string) (domain.PullRequest, string, error) {
	pr, err := s.prRepo.Get(prID)
	if err != nil {
		return domain.PullRequest{}, "", ErrPRNotFound
	}

	if pr.Status == domain.PRStatusMerged {
		return domain.PullRequest{}, "", ErrPRMerged
	}

	if (len(pr.Reviewers) == 0 || pr.Reviewers[0] != oldReviewerID) && (len(pr.Reviewers) < 2  || pr.Reviewers[1] != oldReviewerID) {
		return domain.PullRequest{}, "", ErrReviewerNotAssigned
	}

	oldUser, err := s.userRepo.GetUser(oldReviewerID)
	if err != nil {
		return domain.PullRequest{}, "", ErrReviewerNotAssigned
	}

	excludeIDs := pr.Reviewers
	excludeIDs = append(excludeIDs, pr.AuthorID) 

	teamService := NewTeamService(s.userRepo)
	candidates, _ := teamService.GetRandomActiveMembers(oldUser.TeamName, excludeIDs, 1)
	if len(candidates) == 0 {
		return domain.PullRequest{}, "", ErrNoCandidate
	}

	newReviewerID := candidates[0].ID

	if len(pr.Reviewers) != 0 && pr.Reviewers[0] == oldReviewerID {
		pr.Reviewers[0] = newReviewerID
	} else if len(pr.Reviewers) > 1 && pr.Reviewers[1] == oldReviewerID {
		pr.Reviewers[1] = newReviewerID
	}

	if err := s.prRepo.Update(pr); err != nil {
		return domain.PullRequest{}, "", fmt.Errorf("failed to update PR: %w", err)
	}

	return pr, newReviewerID, nil
}
