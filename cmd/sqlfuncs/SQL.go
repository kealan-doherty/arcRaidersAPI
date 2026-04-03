// file will contain all used public used SQL functions within the API Private SQL functions will be hosted wihtin Pirvate Repo.

package sqlfuncs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ConnectToDB() (*pgx.Conn, error) {
	password := "REDACTED_SECRET"

	dsn := fmt.Sprintf("postgres://%s:%s@REDACTED_HOST:5432/postgres?sslmode=verify-full&sslrootcert=/certs/global-bundle.pem", "REDACTED_USER", password)

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
