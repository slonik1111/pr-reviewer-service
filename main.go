package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/slonik1111/pr-reviewer-service/internal/http/handlers"
	"github.com/slonik1111/pr-reviewer-service/internal/repository/postgres"
	"github.com/slonik1111/pr-reviewer-service/internal/service"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	userRepo := postgres.NewUserRepo(db)
	prRepo := postgres.NewPullRequestRepo(db)

	userSvc := service.NewUserService(userRepo, prRepo)
	teamSvc := service.NewTeamService(userRepo)
	prSvc := service.NewPRService(prRepo, userRepo)

	userHandler := handlers.NewUserHandler(userSvc)
	teamHandler := handlers.NewTeamHandler(teamSvc)
	prHandler := handlers.NewPRHandler(prSvc)

	handlers.RegisterAllHandlers(userHandler, teamHandler, prHandler)

	log.Println("Server running at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
