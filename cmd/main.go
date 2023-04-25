package main

import (
	"github.com/digilob/taskboard-api/pkg/api"
	r "github.com/digilob/taskboard-api/pkg/repository"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {

	// ===================================
	// Connect to the database
	// ===================================

	// Connect to the PostgreSQL database

	postgresDB := ConnectToDB()
	if migErr := HandleMigrations(postgresDB); migErr != nil {
		log.Error("migration error: ", migErr)
	}

	repo := r.NewPostgresTaskRepository(postgresDB)

	// Initialize a panic watcher
	defer recover()

	// Init Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Set up API routes
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/tasks", func(w http.ResponseWriter, r *http.Request) {
			api.GetTasksHandler(repo, w, r)
		})
		apiGroup.POST("/tasks", func(w http.ResponseWriter, r *http.Request) {
			api.CreateTaskHandler(repo, w, r)
		})
		apiGroup.PUT("/tasks/:id", func(w http.ResponseWriter, r *http.Request) {
			api.UpdateTaskHandler(repo, w, r)
		})
		apiGroup.DELETE("/tasks/:id", func(w http.ResponseWriter, r *http.Request) {
			api.DeleteTaskHandler(repo, w, r)
		})
	}

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
