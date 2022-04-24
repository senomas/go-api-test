package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/senomas/go-api/models"
)

type CreateBookInput struct {
	Title   string `json:"title" binding:"required"`
	Author  string `json:"author" binding:"required"`
	Summary string `json:"summary"`
}

type UpdateBookInput struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Summary string `json:"summary"`
}

// GET /books
// POST /books
// Find books
func FindBooks(c *gin.Context) {
	tx := models.DB
	// Validate input
	var query Query
	if str := c.Query("offset"); str != "" {
		if i, err := strconv.Atoi(str); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Offset error": err.Error()})
		} else {
			tx = tx.Offset(i)
		}
	}
	if str := c.Query("limit"); str != "" {
		if i, err := strconv.Atoi(str); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Limit error": err.Error()})
		} else {
			tx = tx.Limit(i)
		}
	} else {
		tx = tx.Limit(1000)
	}
	if c.Request.Method == "POST" {
		if err := c.ShouldBindJSON(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if query.Select != nil {
			tx = tx.Select(query.Select)
		}
		where, params := query.Condition.Apply("", []any{})
		tx.Where(where, params...)
	}

	var books []models.Book
	if err := tx.Find(&books).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": books})
}

// GET /books/:id
// Find a book
func FindBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// PUT /books
// Create new book
func CreateBook(c *gin.Context) {
	// Validate input
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	book := models.Book{Title: input.Title, Author: input.Author, Summary: input.Summary}
	if err := models.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// PATCH /books/:id
// Update a book
func UpdateBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	var input UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.Title = input.Title
	book.Author = input.Author
	if err := models.DB.Updates(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// DELETE /books/:id
// Delete a book
func DeleteBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
