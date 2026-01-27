package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create the Gin Router
	router := gin.Default()

	//Register Routes
	router.GET("/", HomePage)

	//start the API
	router.Run()
	fmt.Printf("server is up and running")

}

func HomePage(c *gin.Context) {
	c.String(http.StatusOK, "this is the home page of the ARC Raiders API router")
}
