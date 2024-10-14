package handlers

import (
	"book-reading-tracker/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloWorld(c *gin.Context) {
	utils.LogInfo("HelloWorld handler called", map[string]interface{}{
		"endpoint": "/hello",
	})

	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
