package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func corsMiddleware(c *gin.Context) {
	allowedOrigins := []string{
		"https://alphawave.gasstrem.com",
		"http://localhost:3000",
		"http://localhost:3001",
	}

	requestOrigin := c.Request.Header.Get("Origin")
	for _, origin := range allowedOrigins {
		if origin == requestOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
			break
		}
	}

	// c.Header("Access-Control-Allow-Origin", "")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
