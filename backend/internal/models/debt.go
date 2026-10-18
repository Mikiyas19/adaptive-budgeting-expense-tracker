package models

import (
	"time"
)

type Debt struct {
	DebtID           int        `json:"debt_id" db:"debt_id"`
	UserID           int        `json:"user_id" db:"user_id"`
	CounterpartyName string     `json:"counterparty_name" db:"counterparty_name" binding:"required"`
	Amount           float64    `json:"amount" db:"amount" binding:"required,gt=0"`
	IsLent           bool       `json:"is_lent" db:"is_lent"` // true = others owe me, false = I owe others
	RemainingAmount  float64    `json:"remaining_amount" db:"remaining_amount"`
	DueDate          *time.Time `json:"due_date" db:"due_date"`
	Status           string     `json:"status" db:"status"` // pending, partial, settled
	Notes            string     `json:"notes" db:"notes"`
	IsSynced         bool       `json:"is_synced" db:"is_synced"`
}
