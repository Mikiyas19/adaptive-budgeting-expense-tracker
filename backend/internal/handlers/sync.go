package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SyncHandler struct {
	DB *sql.DB
}

func NewSyncHandler(db *sql.DB) *SyncHandler {
	return &SyncHandler{DB: db}
}

// GetSyncStatus retrieves the last sync timestamp for a user
func (h *SyncHandler) GetSyncStatus(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	var lastSync sql.NullTime
	err := h.DB.QueryRow(`SELECT last_sync_timestamp FROM users WHERE user_id = $1`, userID).Scan(&lastSync)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"last_sync_timestamp": lastSync.Time,
	})
}

// UploadSync receives local offline creations and stores them in PostgreSQL
func (h *SyncHandler) UploadSync(c *gin.Context) {
	var input struct {
		UserID   int                      `json:"user_id" binding:"required"`
		Expenses []map[string]interface{} `json:"expenses"`
		Debts    []map[string]interface{} `json:"debts"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, _ = h.DB.Exec(`UPDATE users SET last_sync_timestamp = $1 WHERE user_id = $2`, time.Now(), input.UserID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Sync upload completed successfully",
		"status":  "success",
	})
}

// DownloadSync fetches all server records for a user during sync
func (h *SyncHandler) DownloadSync(c *gin.Context) {
	userID := c.Query("user_id")
	lastSyncStr := c.Query("last_sync")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	var lastSyncTime time.Time
	var err error
	if lastSyncStr != "" {
		lastSyncTime, err = time.Parse(time.RFC3339, lastSyncStr)
		if err != nil {
			lastSyncTime = time.Unix(0, 0)
		}
	} else {
		lastSyncTime = time.Unix(0, 0)
	}

	expenseRows, err := h.DB.Query(`SELECT expense_id, user_id, category_id, amount, expense_date, description, is_synced, sync_status 
	                                FROM expenses WHERE user_id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch remote expenses"})
		return
	}
	defer expenseRows.Close()

	var expenses []map[string]interface{}
	for expenseRows.Next() {
		var id, uID, catID int
		var amount float64
		var expDate time.Time
		var desc string
		var isSynced bool
		var status string
		if err := expenseRows.Scan(&id, &uID, &catID, &amount, &expDate, &desc, &isSynced, &status); err == nil {
			expenses = append(expenses, map[string]interface{}{
				"expense_id":   id,
				"user_id":      uID,
				"category_id":  catID,
				"amount":       amount,
				"expense_date": expDate,
				"description":  desc,
				"is_synced":    isSynced,
				"sync_status":  status,
			})
		}
	}

	debtRows, err := h.DB.Query(`SELECT debt_id, user_id, counterparty_name, amount, is_lent, remaining_amount, due_date, status, notes, is_synced 
	                             FROM debts WHERE user_id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch remote debts"})
		return
	}
	defer debtRows.Close()

	var debts []map[string]interface{}
	for debtRows.Next() {
		var id, uID int
		var counterparty string
		var amount, remaining float64
		var isLent, isSynced bool
		var dueDate sql.NullTime
		var status, notes string
		if err := debtRows.Scan(&id, &uID, &counterparty, &amount, &isLent, &remaining, &dueDate, &status, &notes, &isSynced); err == nil {
			debts = append(debts, map[string]interface{}{
				"debt_id":           id,
				"user_id":           uID,
				"counterparty_name": counterparty,
				"amount":            amount,
				"is_lent":           isLent,
				"remaining_amount":  remaining,
				"due_date":          dueDate.Time,
				"status":            status,
				"notes":             notes,
				"is_synced":         isSynced,
			})
		}
	}

	_ = lastSyncTime // Keep reference if needed for filtering

	c.JSON(http.StatusOK, gin.H{
		"server_time": time.Now(),
		"expenses":    expenses,
		"debts":       debts,
	})
}
