package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	"github.com/slonik1111/pr-reviewer-service/internal/service"
)

type PRHandler struct {
	prService *service.PRService
}

func NewPRHandler(prSvc *service.PRService) *PRHandler {
	return &PRHandler{prService: prSvc}
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	
	pr := domain.PullRequest{
		ID:   req.PullRequestID,
		Title: req.PullRequestName,
		AuthorID: req.AuthorID,
	}

	pr, err := h.prService.CreatePR(pr)

	if  err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	ret := retPR{
		ID: pr.ID,
		Description: pr.Description,	
		AuthorID: pr.AuthorID,
		Status: string(pr.Status),
		Reviewers: pr.Reviewers,
		MergedAt: nil,
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"pr": ret})
}

func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pr, err := h.prService.MergePR(req.PullRequestID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}	

	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": retPR{
			ID: pr.ID,
			Description: pr.Description,
			AuthorID: pr.AuthorID,
			Status: string(pr.Status),
			Reviewers: pr.Reviewers,
			MergedAt: pr.MergedAt,
		},
	})
}

func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldReviewerID string `json:"old_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pr, newReviewerID, err := h.prService.ReassignReviewer(req.PullRequestID, req.OldReviewerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	ret := retPR{
		ID: pr.ID,
		Description: pr.Description,	
		AuthorID: pr.AuthorID,
		Status: string(pr.Status),
		Reviewers: pr.Reviewers,
		MergedAt: nil,
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr":          ret,
		"replaced_by": newReviewerID,
	})
}
