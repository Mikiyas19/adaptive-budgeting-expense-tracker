package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Mikiyas19/adaptive-budgeting-expense-tracker/backend/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

// Register handles new user sign-up (FR-01, FR-02, FR-03)
func (h *AuthHandler) Register(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	query := `INSERT INTO users (email, password_hash, full_name, created_at) VALUES ($1, $2, $3, $4) RETURNING user_id`
	var userID int
	err = h.DB.QueryRow(query, input.Email, string(hashedPassword), input.FullName, time.Now()).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered or database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": userID,
	})
}

// Login handles user authentication (FR-04)
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedUser models.User
	query := `SELECT user_id, email, password_hash, full_name FROM users WHERE email = $1`
	err := h.DB.QueryRow(query, input.Email).Scan(&storedUser.UserID, &storedUser.Email, &storedUser.PasswordHash, &storedUser.FullName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Login successful",
		"user_id":   storedUser.UserID,
		"email":     storedUser.Email,
		"full_name": storedUser.FullName,
	})
}

// ForgotPassword handles requesting a password reset token (FR-05)
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resetToken := "123456" // Development placeholder token
	expiresAt := time.Now().Add(15 * time.Minute)

	var userID int
	err := h.DB.QueryRow(`SELECT user_id FROM users WHERE email = $1`, input.Email).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset code has been sent."})
		return
	}

	_, err = h.DB.Exec(`UPDATE users SET reset_token = $1, reset_expires = $2 WHERE user_id = $3`, resetToken, expiresAt, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password reset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Password reset code generated successfully",
		"reset_token": resetToken,
	})
}

// ResetPassword handles updating the password using the reset token
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input struct {
		Email       string `json:"email" binding:"required,email"`
		ResetToken  string `json:"reset_token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID int
	var storedToken string
	var expiresAt time.Time

	err := h.DB.QueryRow(`SELECT user_id, reset_token, reset_expires FROM users WHERE email = $1`, input.Email).Scan(&userID, &storedToken, &expiresAt)
	if err != nil || storedToken != input.ResetToken || time.Now().After(expiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	_, err = h.DB.Exec(`UPDATE users SET password_hash = $1, reset_token = NULL, reset_expires = NULL WHERE user_id = $2`, string(hashedPassword), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
