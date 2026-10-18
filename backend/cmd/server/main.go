package main

import (
	"log"
	"net/http"

	"github.com/Mikiyas19/adaptive-budgeting-expense-tracker/backend/internal/config"
	"github.com/Mikiyas19/adaptive-budgeting-expense-tracker/backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to PostgreSQL
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)

	// API Version 1 Group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "API is running"})
		})

		// Auth routes
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
	}

	// Start server
	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}