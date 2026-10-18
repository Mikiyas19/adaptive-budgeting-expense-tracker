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

	// Hash the password using bcrypt (NFR-06)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user into database
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

	// Verify password hash
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