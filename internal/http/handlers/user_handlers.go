package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/slonik1111/pr-reviewer-service/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userService: userSvc}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.userService.SetIsActive(req.UserID, req.IsActive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	req.IsActive = user.IsActive
	json.NewEncoder(w).Encode(req)
}

type retPR struct {
	ID          string `json:"pull_request_id"`
	Description string `json:"pull_request_name"`
	AuthorID    string `json:"author_id"`
	Status      string `json:"status"`
	Reviewers	[]string `json:"revievers"`
	MergedAt    *string  `json:"mergedAt,omitempty"`
}

func (h *UserHandler) GetPRsForUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}
	prs, err := h.userService.GetPRsForReview(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	retPrs := make([]retPR, len(prs))
	for i := range prs {
		retPrs[i] = retPR{
			ID: prs[i].ID,
			Description: prs[i].Description,
			AuthorID: prs[i].AuthorID,
			Status: string(prs[i].Status),
			Reviewers: prs[i].Reviewers,
			MergedAt: nil,
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":       userID,
		"pull_requests": retPrs,
	})
}
