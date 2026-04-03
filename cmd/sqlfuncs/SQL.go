// file will contain all used public used SQL functions within the API Private SQL functions will be hosted wihtin Pirvate Repo.

package sqlfuncs

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectToDB() (*pgx.Conn, error) {
	DB_CONNECT := os.Getenv("DB_CONNECT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")

	if DB_CONNECT == "" || DB_USER == "" || DB_PASSWORD == "" {
		return nil, fmt.Errorf("database env vars DB_CONNECT, DB_USER, and DB_PASSWORD are required")
	}

	dsn := fmt.Sprintf(DB_CONNECT, DB_USER, DB_PASSWORD)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func DisconnectDB(conn *pgx.Conn) error {
	err := conn.Close(context.Background())
	if err != nil {
		return err
	}
	return nil
}
