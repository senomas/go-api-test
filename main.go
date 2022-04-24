package main

import (
	"github.com/gin-gonic/gin"
	"github.com/senomas/go-api/controllers"
	"github.com/senomas/go-api/models"
)

func xmain() {
	r := gin.Default()

	models.ConnectDatabase()
	controllers.SetupRoutes(r)

	r.Run()
}
