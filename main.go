package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"zerodependency.co.uk/spectre/pkg/spectre"
)

const (
	ipAddress = "127.0.0.1"
	port      = 18080
)

func engine() *gin.Engine {
	return gin.Default()
}

func main() {
	spectre.Init("definitions/")

	server := engine()

	v1api := server.Group("/api/v1/spectre")
	v1api.GET("/tests/:service", httpGetSpectreTestsForService)
	v1api.POST("/tests/:id/invoke", httpPostInvokeSpectreTest)

	server.Run(fmt.Sprintf("%v:%v", ipAddress, port))
}
