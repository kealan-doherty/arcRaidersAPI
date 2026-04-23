// file will contain all used public used SQL functions within the API Private SQL functions will be hosted wihtin Pirvate Repo.

package sqlfuncs

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

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

func CreateTable(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS items (
			id TEXT PRIMARY KEY,
			type TEXT,
			rarity TEXT,
			weightkg NUMERIC(5,2),
			value NUMERIC,
			isweapon BOOLEAN
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

func AddItems(conn *pgx.Conn) error {
	csvPath := strings.TrimSpace(os.Getenv("ITEMS_CSV_PATH"))
	if csvPath == "" {
		csvPath = "data.csv"
	}

	if cwd, err := os.Getwd(); err == nil {
		log.Printf("AddItems: cwd=%s csv_path=%s", cwd, csvPath)
	}

	file, err := os.Open(csvPath)
	if err != nil && csvPath == "data.csv" {
		fallback := "/data.csv"
		log.Printf("AddItems: primary csv path unavailable, trying fallback path=%s", fallback)
		file, err = os.Open(fallback)
		if err == nil {
			csvPath = fallback
		}
	}
	if err != nil {
		return fmt.Errorf("open data.csv at %s: %w", csvPath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("read csv header: %w", err)
	}

	rows := make([][]any, 0, 1024)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read csv row: %w", err)
		}

		if len(record) != 6 {
			return fmt.Errorf("expected 6 columns, got %d", len(record))
		}

		weightKg, err := parseNullableFloat(record[3])
		if err != nil {
			return fmt.Errorf("invalid weightkg for item %q: %w", record[0], err)
		}

		value, err := parseNullableFloat(record[4])
		if err != nil {
			return fmt.Errorf("invalid value for item %q: %w", record[0], err)
		}

		isWeapon, err := strconv.ParseBool(strings.TrimSpace(record[5]))
		if err != nil {
			return fmt.Errorf("invalid isWeapon for item %q: %w", record[0], err)
		}

		rows = append(rows, []any{
			strings.TrimSpace(record[0]),
			strings.TrimSpace(record[1]),
			strings.TrimSpace(record[2]),
			weightKg,
			value,
			isWeapon,
		})
	}

	if len(rows) == 0 {
		log.Printf("AddItems: no rows parsed from csv")
		return nil
	}

	log.Printf("AddItems: parsed rows=%d", len(rows))

	copiedCount, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"items"},
		[]string{"id", "type", "rarity", "weightkg", "value", "isWeapon"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("copy items into table: %w", err)
	}

	log.Printf("AddItems: copied rows=%d", copiedCount)

	return nil
}

func parseNullableFloat(raw string) (any, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
