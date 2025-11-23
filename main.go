package main

import (
	"log"
	"net/http"

	"github.com/slonik1111/pr-reviewer-service/internal/http/handlers"
	"github.com/slonik1111/pr-reviewer-service/internal/repository/inmemory"
	"github.com/slonik1111/pr-reviewer-service/internal/service"
)

func main() {
	userRepo := inmemory.NewUserRepoInMemory()
	teamRepo := inmemory.NewTeamRepoInMemory()
	prRepo := inmemory.NewPRRepoInMemory()

	userSvc := service.NewUserService(userRepo, prRepo)
	teamSvc := service.NewTeamService(teamRepo, userRepo)
	prSvc := service.NewPRService(prRepo, userRepo, teamRepo)

	userHandler := handlers.NewUserHandler(userSvc)
	teamHandler := handlers.NewTeamHandler(teamSvc)
	prHandler := handlers.NewPRHandler(prSvc)

	handlers.RegisterAllHandlers(userHandler, teamHandler, prHandler)

	log.Println("Server running at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
