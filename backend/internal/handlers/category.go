package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	DB *sql.DB
}

func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

// CreateCategory handles creating a new personal category (FR-43)
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input struct {
		UserID  int    `json:"user_id" binding:"required"`
		Name    string `json:"name" binding:"required"`
		Icon    string `json:"icon"`
		IsFixed bool   `json:"is_fixed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO categories (user_id, name, icon, is_predefined, is_fixed) VALUES ($1, $2, $3, FALSE, $4) RETURNING category_id`
	var categoryID int
	err := h.DB.QueryRow(query, input.UserID, input.Name, input.Icon, input.IsFixed).Scan(&categoryID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Category already exists or database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Category created successfully",
		"category_id": categoryID,
	})
}

// GetCategories returns predefined categories and user-specific custom categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID := c.Query("user_id")

	query := `SELECT category_id, user_id, name, icon, is_predefined, is_fixed FROM categories WHERE is_predefined = TRUE OR user_id = $1`
	rows, err := h.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}
	defer rows.Close()

	var categories []map[string]interface{}
	for rows.Next() {
		var catID int
		var uID sql.NullInt64
		var name, icon string
		var isPredefined, isFixed bool

		if err := rows.Scan(&catID, &uID, &name, &icon, &isPredefined, &isFixed); err != nil {
			continue
		}

		categories = append(categories, map[string]interface{}{
			"category_id":   catID,
			"user_id":       uID.Int64,
			"name":          name,
			"icon":          icon,
			"is_predefined": isPredefined,
			"is_fixed":      isFixed,
		})
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
