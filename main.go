package main

import (
	"github.com/senomas/go-api/controllers"
	"github.com/senomas/go-api/models"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/books", controllers.FindBooks)
	r.POST("/books", controllers.FindBooks)
	r.GET("/books/:id", controllers.FindBook)
	r.PUT("/books", controllers.CreateBook)
	r.PATCH("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)
}

func main() {
	r := gin.Default()

	models.ConnectDatabase()

	SetupRoutes(r)

	r.Run()
}
