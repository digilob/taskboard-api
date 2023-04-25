package main

import (
	"github.com/digilob/taskboard-api/pkg/api"
	"github.com/gin-gonic/gin"
)

func main() {

	// ===================================
	// Connect to the database
	// ===================================

	// Connect to the PostgreSQL database
	postgresDB := connectToDB()
	if migErr := handleMigrations(postgresDB); migErr != nil {
		log.Error("migration error: ", migErr)
	}

	repo := NewRepository(postgresDB)

	// Initialize a panic watcher
	defer recover()

	// Init Gin router
	r := gin.Default()

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/tasks", api.GetTasksHandler)
		apiGroup.POST("/tasks", api.CreateTaskHandler)
		apiGroup.PUT("/tasks/:id", api.UpdateTaskHandler)
		apiGroup.DELETE("/tasks/:id", api.DeleteTaskHandler)
	}

	r.Run(":8080")
}
