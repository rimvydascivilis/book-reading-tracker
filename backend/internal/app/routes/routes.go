package routes

import (
	"book-reading-tracker/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/hello", handlers.HelloWorld)
}
