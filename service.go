package main

import (
	"net/http"

	"github.com/ZeroDependency/spectre/pkg/spectre"
	"github.com/gin-gonic/gin"
)

func httpGetSpectreTestsForService(c *gin.Context) {
	service := c.Param("service")
	if service == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tests := spectre.GetSpectreTestsForService(service)

	if len(tests) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, tests)
}

func httpPostInvokeSpectreTest(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := spectre.InvokeSpectreTest(ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
