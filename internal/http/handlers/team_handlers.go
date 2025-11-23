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

func NewTeamHandler(teamSvc *service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamSvc}
}

type reqUser struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type reqTeam struct {
	Name 	string  `json:"team_name"`
	Users 	[]reqUser	`json:"members"`
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

	team := req.TeamName
	members := make([]domain.User, len(req.Members))
	for i := range req.Members {
		members[i] = domain.User{
			ID:       req.Members[i].ID,
			Name:     req.Members[i].Name,
			TeamName: req.TeamName,
			IsActive: req.Members[i].IsActive}
	}

	if err := h.teamService.CreateTeam(team, members); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	retTeam := reqTeam{Name: team, Users: req.Members}

	json.NewEncoder(w).Encode(retTeam)
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
	retTeam := reqTeam{
		Name: teamName,
		Users: make([]reqUser, len(team)),
	}
	for i := range team {
		retTeam.Users[i] = reqUser {
			ID: team[i].ID,
			Name: team[i].Name,
			IsActive: team[i].IsActive,
		}
	}
	json.NewEncoder(w).Encode(retTeam)
}
