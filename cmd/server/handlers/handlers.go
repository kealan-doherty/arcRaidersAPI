package handlers

import (
	"arcRaidersAPI/cmd/sqlfuncs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// basic Ping route for fun
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func GetItems(c *gin.Context, conn *pgx.Conn) {
	items, err := sqlfuncs.GetAllItems(conn)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, items)
}
