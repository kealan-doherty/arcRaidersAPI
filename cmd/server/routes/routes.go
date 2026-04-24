package routes

import (
	"arcRaidersAPI/cmd/server/handlers"

	"github.com/gin-gonic/gin"
)

// will handle all routes in the Arc Raiders API

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", handlers.Ping)
	r.GET("/items", handlers.GetItems)
}
