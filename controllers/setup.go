package controllers

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	r.GET("/books", FindBooks)
	r.POST("/books", FindBooks)
	r.GET("/books/:id", FindBook)
	r.PUT("/books", CreateBook)
	r.PATCH("/books/:id", UpdateBook)
	r.DELETE("/books/:id", DeleteBook)
}
