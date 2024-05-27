package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/marcboeker/go-duckdb"
)

func main() {
	ctx := context.Background()
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	currentTimestamp, err := run(ctx)
	if err != nil {
		l.Error("error", "ERROR", err)
		os.Exit(1)
	}
	l.Info("current timestamp", "timestamp", currentTimestamp)
}

func run(ctx context.Context) (currentTimestamp string, ferr error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return "", fmt.Errorf("opening duckdb connection: %w", err)
	}

	defer func(db *sql.DB) {
		if cerr := db.Close(); cerr != nil {
			ferr = errors.Join(ferr, fmt.Errorf("closing duckdb connection: %w", cerr))
		}
	}(db)

	row := db.QueryRowContext(ctx, "SELECT current_timestamp")

	if err = row.Scan(&currentTimestamp); err != nil {
		return "", fmt.Errorf("scanning row: %w", err)
	}

	return currentTimestamp, nil
}
