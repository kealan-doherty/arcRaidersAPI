package main

import (
	"arcRaidersAPI/cmd/sqlfuncs"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	conn, err := sqlfuncs.ConnectToDB()

	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer func() {
		if err := sqlfuncs.DisconnectDB(conn); err != nil {
			log.Printf("Unable to disconnect from database: %v", err)
		}
	}()

	sqlfuncs.AddItems(conn)

	r.Run()

}
