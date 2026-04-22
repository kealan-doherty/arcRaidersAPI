package main

import (
	"arcRaidersAPI/cmd/sqlfuncs"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	conn, err := sqlfuncs.ConnectToDB()

	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer func() {
		if err := sqlfuncs.DisconnectDB(conn); err != nil {
			log.Printf("Unable to disconnect from database: %v", err)
		}
	}()

	if err := sqlfuncs.CreateTable(conn); err != nil {
		log.Fatalf("Create table failed: %v", err)
	}

	var v string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&v); err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Println(v)

	r.Run()

}
