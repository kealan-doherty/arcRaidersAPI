package routes

import (
	"arcRaidersAPI/cmd/server/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// will handle all routes in the Arc Raiders API

func RegisterRoutes(r *gin.Engine, conn *pgx.Conn) {
	r.GET("/ping", handlers.Ping)
	r.GET("/items", func(c *gin.Context) {
		handlers.GetItems(c, conn)
	})
}
