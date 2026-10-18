package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Mikiyas19/adaptive-budgeting-expense-tracker/backend/internal/models"
	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	DB *sql.DB
}

func NewExpenseHandler(db *sql.DB) *ExpenseHandler {
	return &ExpenseHandler{DB: db}
}

// AddExpense handles recording a new expense (FR-08)
func (h *ExpenseHandler) AddExpense(c *gin.Context) {
	var input models.Expense
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		formatDate := time.Now()
		if input.ExpenseDate.IsZero() {
			input.ExpenseDate = formatDate
		}
		_ = formatDate // placeholder to use variable if needed
	}

	if input.ExpenseDate.IsZero() {
		input.ExpenseDate = time.Now()
	}

	query := `INSERT INTO expenses (user_id, category_id, amount, expense_date, description, is_synced, sync_status) 
	          VALUES ($1, $2, $3, $4, $5, TRUE, 'synced') RETURNING expense_id`
	var expenseID int
	err := h.DB.QueryRow(query, input.UserID, input.CategoryID, input.Amount, input.ExpenseDate, input.Description).Scan(&expenseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save expense: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Expense added successfully",
		"expense_id": expenseID,
	})
}

// GetExpenses handles retrieving all expenses for a user (FR-12)
func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	rows, err := h.DB.Query(`SELECT expense_id, user_id, category_id, amount, expense_date, description, is_synced, sync_status 
	                         FROM expenses WHERE user_id = $1 ORDER BY expense_date DESC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses"})
		return
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ExpenseID, &e.UserID, &e.CategoryID, &e.Amount, &e.ExpenseDate, &e.Description, &e.IsSynced, &e.SyncStatus); err != nil {
			continue
		}
		expenses = append(expenses, e)
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

// DeleteExpense handles deleting an expense (FR-11)
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	expenseID := c.Param("id")
	_, err := h.DB.Exec(`DELETE FROM expenses WHERE expense_id = $1`, expenseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
