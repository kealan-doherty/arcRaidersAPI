package handlers

import (
	"github.com/gin-gonic/gin"
)

// basic Ping route for fun
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func GetItems(c *gin.Context) {

}
