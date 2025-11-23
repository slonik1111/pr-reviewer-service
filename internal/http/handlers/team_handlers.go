package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/slonik1111/pr-reviewer-service/internal/domain"
	"github.com/slonik1111/pr-reviewer-service/internal/service"
)

type TeamHandler struct {
	teamService *service.TeamService
}

type reqUser struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func NewTeamHandler(teamSvc *service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamSvc}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TeamName string    `json:"team_name"`
		Members  []reqUser `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	team := &domain.Team{Name: req.TeamName}
	members := make([]*domain.User, len(req.Members))
	for i := range req.Members {
		members[i] = &domain.User{
			ID:       req.Members[i].ID,
			Name:     req.TeamName,
			TeamID:   req.Members[i].Name,
			IsActive: req.Members[i].IsActive}
	}

	if err := h.teamService.CreateTeam(team, members); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"team": team})
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		http.Error(w, "team_name required", http.StatusBadRequest)
		return
	}
	team, err := h.teamService.GetTeam(teamName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(team)
}
