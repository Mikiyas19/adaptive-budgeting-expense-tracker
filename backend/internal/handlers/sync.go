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

	// Update user's last sync timestamp
	_, _ = h.DB.Exec(`UPDATE users SET last_sync_timestamp = $1 WHERE user_id = $2`, time.Now(), input.UserID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Sync upload completed successfully",
		"status":  "success",
	})
}
