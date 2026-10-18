package models

import (
	"time"
)

type User struct {
	UserID             int       `json:"user_id" db:"user_id"`
	Email              string    `json:"email" db:"email" binding:"required,email"`
	PasswordHash       string    `json:"-" db:"password_hash"`
	Password           string    `json:"password,omitempty" binding:"required,min=6"` // Used only during registration/login payload
	FullName           string    `json:"full_name" db:"full_name"`
	MonthlyBudget      float64   `json:"monthly_budget" db:"monthly_budget"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	LastSyncTimestamp  *time.Time `json:"last_sync_timestamp" db:"last_sync_timestamp"`
}