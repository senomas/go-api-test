package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/senomas/go-api/models"
	"gorm.io/gorm"
)

type FindQuery struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit" binding:"required"`
}

type QueryBookInput struct {
	FindQuery
	Title_Like  string `json:"title-like"`
	Author_Like string `json:"author-like"`
}

type CreateBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

// POST /books
// Find books
func FindBooks(c *gin.Context) {
	var books []models.Book
	var tx *gorm.DB
	// Validate input
	var input QueryBookInput
	if c.Request.Method == "POST" {
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tx = models.DB.Offset(input.Offset).Limit(input.Limit)
		query := strings.Builder{}
		var params []any
		if input.Title_Like != "" {
			if query.Len() > 0 {
				query.WriteString(" AND ")
			}
			query.WriteString("title like ?")
			params = append(params, input.Title_Like)
		}
		if input.Author_Like != "" {
			if query.Len() > 0 {
				query.WriteString(" AND ")
			}
			query.WriteString("author like ?")
			params = append(params, input.Author_Like)
		}
		tx.Where(query.String(), params...)
	} else {
		tx = models.DB.Limit(1000)
	}

	tx.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}

// GET /books/:id
// Find a book
func FindBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
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
	book := models.Book{Title: input.Title, Author: input.Author}
	models.DB.Create(&book)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// PATCH /books/:id
// Update a book
func UpdateBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
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
	models.DB.Updates(&book)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// DELETE /books/:id
// Delete a book
func DeleteBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	models.DB.Delete(&book)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
