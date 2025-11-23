package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	"github.com/slonik1111/pr-reviewer-service/internal/repository"
)

var (
	ErrPRNotFound      = errors.New("pull request not found")
	ErrPRMerged        = errors.New("cannot modify a merged PR")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate     = errors.New("no active replacement candidate in team")
)

// PRService управляет Pull Request и назначением ревьюверов
type PRService struct {
	prRepo    repo.PullRequestRepository
	userRepo  repo.UserRepository
	teamRepo  repo.TeamRepository
}

// NewPRService создает сервис PR
func NewPRService(prRepo repo.PullRequestRepository, userRepo repo.UserRepository, teamRepo repo.TeamRepository) *PRService {
	return &PRService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// CreatePR создает новый Pull Request и назначает до 2 ревьюверов
func (s *PRService) CreatePR(pr *domain.PullRequest) error {
	// проверка существования PR
	_, err := s.prRepo.Get(pr.ID)
	if err == nil {
		return fmt.Errorf("PR %s already exists", pr.ID)
	}

	// проверка существования автора
	author, err := s.userRepo.GetUser(pr.AuthorID)
	if err != nil {
		return fmt.Errorf("author %s not found", pr.AuthorID)
	}

	pr.TeamID = author.TeamID

	// получаем до 2 ревьюверов
	teamService := NewTeamService(s.teamRepo, s.userRepo)
	reviewers, _ := teamService.GetRandomActiveMembers(pr.TeamID, []string{pr.AuthorID}, 2)

	if len(reviewers) > 0 {
		pr.Reviewer1 = &reviewers[0].ID
	}
	if len(reviewers) > 1 {
		pr.Reviewer2 = &reviewers[1].ID
	}

	pr.Status = domain.PRStatusOpen
	log.Println("revievers: ", reviewers)
	pr.CreatedAt = time.Now().Unix()

	return s.prRepo.Create(pr)
}

// MergePR помечает PR как MERGED (идемпотентно)
func (s *PRService) MergePR(prID string) (*domain.PullRequest, error) {
	pr, err := s.prRepo.Get(prID)
	if err != nil {
		return nil, ErrPRNotFound
	}

	if pr.Status == domain.PRStatusMerged {
		// идемпотентно
		return pr, nil
	}

	now := time.Now().Unix()
	pr.Status = domain.PRStatusMerged
	pr.MergedAt = &now

	if err := s.prRepo.Update(pr); err != nil {
		return nil, fmt.Errorf("failed to merge PR: %w", err)
	}

	return pr, nil
}

// ReassignReviewer заменяет одного ревьювера на другого из команды
func (s *PRService) ReassignReviewer(prID string, oldReviewerID string) (*domain.PullRequest, string, error) {
	pr, err := s.prRepo.Get(prID)
	if err != nil {
		return nil, "", ErrPRNotFound
	}

	if pr.Status == domain.PRStatusMerged {
		return nil, "", ErrPRMerged
	}

	// проверяем, что oldReviewer действительно назначен
	if (pr.Reviewer1 == nil || *pr.Reviewer1 != oldReviewerID) && (pr.Reviewer2 == nil || *pr.Reviewer2 != oldReviewerID) {
		return nil, "", ErrReviewerNotAssigned
	}

	// получаем кандидатов для замены
	oldUser, err := s.userRepo.GetUser(oldReviewerID)
	if err != nil {
		return nil, "", ErrReviewerNotAssigned
	}

	teamService := NewTeamService(s.teamRepo, s.userRepo)
	candidates, _ := teamService.GetRandomActiveMembers(oldUser.TeamID, []string{oldReviewerID}, 1)
	if len(candidates) == 0 {
		return nil, "", ErrNoCandidate
	}

	newReviewerID := candidates[0].ID

	// меняем ревьювера
	if pr.Reviewer1 != nil && *pr.Reviewer1 == oldReviewerID {
		pr.Reviewer1 = &newReviewerID
	} else if pr.Reviewer2 != nil && *pr.Reviewer2 == oldReviewerID {
		pr.Reviewer2 = &newReviewerID
	}

	if err := s.prRepo.Update(pr); err != nil {
		return nil, "", fmt.Errorf("failed to update PR: %w", err)
	}

	return pr, newReviewerID, nil
}
