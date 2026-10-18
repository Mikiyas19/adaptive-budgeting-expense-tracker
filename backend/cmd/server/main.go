package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize Gin router
    r := gin.Default()

    // Health check endpoint
    r.GET("/api/v1/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":  "success",
            "message": "Adaptive Budgeting API is running smoothly",
        })
    })

    // Start server on port 8080
    log.Println("Server starting on port 8080...")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}