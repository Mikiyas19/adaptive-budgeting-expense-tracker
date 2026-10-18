package main

import (
	"log"
	"net/http"

	"github.com/Mikiyas19/adaptive-budgeting-expense-tracker/backend/internal/config"
	"github.com/Mikiyas19/adaptive-budgeting-expense-tracker/backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	authHandler := handlers.NewAuthHandler(db)
	expenseHandler := handlers.NewExpenseHandler(db)
	debtHandler := handlers.NewDebtHandler(db)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "API is running"})
		})

		// Auth routes
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)

		// Expense routes
		v1.POST("/expenses", expenseHandler.AddExpense)
		v1.GET("/expenses", expenseHandler.GetExpenses)
		v1.DELETE("/expenses/:id", expenseHandler.DeleteExpense)

		// Debt routes
		v1.POST("/debts", debtHandler.AddDebt)
		v1.GET("/debts", debtHandler.GetDebts)
		v1.PATCH("/debts/:id/settle", debtHandler.SettleDebt)
	}

	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
