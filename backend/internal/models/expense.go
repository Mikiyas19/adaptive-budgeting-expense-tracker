package models

import (
	"time"
)

type Expense struct {
	ExpenseID   int       `json:"expense_id" db:"expense_id"`
	UserID      int       `json:"user_id" db:"user_id"`
	CategoryID  int       `json:"category_id" db:"category_id" binding:"required"`
	Amount      float64   `json:"amount" db:"amount" binding:"required,gt=0"`
	ExpenseDate time.Time `json:"expense_date" db:"expense_date"`
	Description string    `json:"description" db:"description"`
	IsSynced    bool      `json:"is_synced" db:"is_synced"`
	SyncStatus  string    `json:"sync_status" db:"sync_status"`
}
