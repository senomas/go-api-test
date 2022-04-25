package controllers

import (
	"net/http"

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
	var books []models.Book
	models.DB.Finds(c, &models.Book{}, &books)
}

// GET /books/:id
// Find a book
func FindBook(c *gin.Context) {
	var book models.Book
	models.DB.Find(c, &book)
}

// PUT /books
// Create new book
func CreateBook(c *gin.Context) {
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := models.Book{Title: input.Title, Author: input.Author, Summary: input.Summary}
	models.DB.Create(c, &book)
}

// PATCH /books/:id
// Update a book
func UpdateBook(c *gin.Context) {
	var input UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var book models.Book
	models.DB.Update(c, &book, func() {
		book.Title = input.Title
		book.Author = input.Author
	})
}

// DELETE /books/:id
// Delete a book
func DeleteBook(c *gin.Context) {
	var book models.Book
	models.DB.Delete(c, &book)
}
