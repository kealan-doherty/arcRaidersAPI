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

	if cfg := conn.Config(); cfg != nil {
		log.Printf("DB target host=%s port=%d db=%s user=%s", cfg.Host, cfg.Port, cfg.Database, cfg.User)
	}

	if err := sqlfuncs.AddItems(conn); err != nil {
		log.Fatalf("failed to add items to table: %v", err)
	}

	var itemCount int64
	if err := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM items").Scan(&itemCount); err != nil {
		log.Fatalf("failed to count items: %v", err)
	}
	log.Printf("items row count after AddItems: %d", itemCount)

	var v string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&v); err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Println(v)

	r.Run()

}
