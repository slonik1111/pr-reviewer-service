package handlers

import (
	"encoding/json"
	"log"
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
		UserID  string `json:"user_id"`
		IsActive bool  `json:"is_active"`
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
	json.NewEncoder(w).Encode(map[string]interface{}{"user": user})
}

func (h *UserHandler) GetPRsForUser(w http.ResponseWriter, r *http.Request) {
	log.Println(h.userService.GetUser("u1"))
	log.Println(h.userService.GetUser("u2"))
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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
