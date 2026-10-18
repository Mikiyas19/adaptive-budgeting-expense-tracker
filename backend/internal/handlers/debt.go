package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Mikiyas19/adaptive-budgeting-expense-tracker/backend/internal/models"
	"github.com/gin-gonic/gin"
)

type DebtHandler struct {
	DB *sql.DB
}

func NewDebtHandler(db *sql.DB) *DebtHandler {
	return &DebtHandler{DB: db}
}

// AddDebt handles recording borrowed or lent money (FR-14, FR-15)
func (h *DebtHandler) AddDebt(c *gin.Context) {
	var input models.Debt
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set initial remaining amount equal to original amount
	input.RemainingAmount = input.Amount
	input.Status = "pending"

	query := `INSERT INTO debts (user_id, counterparty_name, amount, is_lent, remaining_amount, due_date, status, notes, is_synced) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE) RETURNING debt_id`
	var debtID int
	err := h.DB.QueryRow(query, input.UserID, input.CounterpartyName, input.Amount, input.IsLent, input.RemainingAmount, input.DueDate, input.Status, input.Notes).Scan(&debtID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record debt: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Debt recorded successfully",
		"debt_id": debtID,
	})
}

// GetDebts handles retrieving all peer debts for a user (FR-16)
func (h *DebtHandler) GetDebts(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	rows, err := h.DB.Query(`SELECT debt_id, user_id, counterparty_name, amount, is_lent, remaining_amount, due_date, status, notes, is_synced 
	                         FROM debts WHERE user_id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve debts"})
		return
	}
	defer rows.Close()

	var debts []models.Debt
	for rows.Next() {
		var d models.Debt
		if err := rows.Scan(&d.DebtID, &d.UserID, &d.CounterpartyName, &d.Amount, &d.IsLent, &d.RemainingAmount, &d.DueDate, &d.Status, &d.Notes, &d.IsSynced); err != nil {
			continue
		}
		debts = append(debts, d)
	}

	c.JSON(http.StatusOK, gin.H{"debts": debts})
}

// SettleDebt marks a debt as settled or processes repayment (FR-17, FR-18)
func (h *DebtHandler) SettleDebt(c *gin.Context) {
	debtID := c.Param("id")

	_, err := h.DB.Exec(`UPDATE debts SET remaining_amount = 0, status = 'settled' WHERE debt_id = $1`, debtID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to settle debt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Debt marked as settled successfully"})
}
