package main

import (
	"github.com/gin-gonic/gin"
	"github.com/senomas/go-api/controllers"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"0.0.0.0"})

	controllers.SetupRoutes(r)

	r.Run()
}
