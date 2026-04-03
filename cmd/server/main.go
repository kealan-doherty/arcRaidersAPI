package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {

	DB_CONNECT := os.Getenv("DB_CONNECT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")

	// For IAM/Secrets: inject via env or fetch before run.
	dsn := fmt.Sprintf(DB_CONNECT, DB_USER, DB_PASSWORD)

	conn, err := pgx.Connect(context.Background(), dsn)
	
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close(context.Background())

	var v string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&v); err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Println(v)

}
