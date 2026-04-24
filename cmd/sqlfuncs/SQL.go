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
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			type TEXT NOT NULL,
			rarity TEXT NOT NULL,
			weightkg NUMERIC(5,2),
			value NUMERIC,
			isweapon BOOLEAN NOT NULL
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
	reader.FieldsPerRecord = -1
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("read csv header: %w", err)
	}

	rows := make([][]any, 0, 1024)
	skippedRows := 0
	lineNum := 1
	for {
		record, err := reader.Read()
		lineNum++
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read csv row: %w", err)
		}

		if len(record) != 6 {
			log.Printf("AddItems: skipping line %d due to invalid column count (got=%d want=6): %v", lineNum, len(record), record)
			skippedRows++
			continue
		}

		weightKg, err := parseNullableFloat(record[3])
		if err != nil {
			log.Printf("AddItems: skipping line %d due to invalid weightkg for item %q: %v", lineNum, record[0], err)
			skippedRows++
			continue
		}

		value, err := parseNullableFloat(record[4])
		if err != nil {
			log.Printf("AddItems: skipping line %d due to invalid value for item %q: %v", lineNum, record[0], err)
			skippedRows++
			continue
		}

		isWeapon, err := strconv.ParseBool(strings.TrimSpace(record[5]))
		if err != nil {
			log.Printf("AddItems: skipping line %d due to invalid isWeapon for item %q: %v", lineNum, record[0], err)
			skippedRows++
			continue
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

	log.Printf("AddItems: parsed valid rows=%d skipped rows=%d", len(rows), skippedRows)

	ctx := context.Background()
	batch := &pgx.Batch{}
	for _, row := range rows {
		batch.Queue(
			`INSERT INTO items (name, type, rarity, weightkg, value, isweapon)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (name)
			 DO UPDATE SET
			 	type = EXCLUDED.type,
			 	rarity = EXCLUDED.rarity,
			 	weightkg = EXCLUDED.weightkg,
			 	value = EXCLUDED.value,
			 	isweapon = EXCLUDED.isweapon`,
			row...,
		)
	}

	batchResults := conn.SendBatch(ctx, batch)
	for i := 0; i < len(rows); i++ {
		if _, err := batchResults.Exec(); err != nil {
			_ = batchResults.Close()
			return fmt.Errorf("upsert item at batch index %d: %w", i, err)
		}
	}
	if err := batchResults.Close(); err != nil {
		return fmt.Errorf("finalize item upsert batch: %w", err)
	}

	log.Printf("AddItems: upserted rows=%d", len(rows))

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
