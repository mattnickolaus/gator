package sql

import (
	"database/sql"
	"embed"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed schema
var embedMigrations embed.FS

func UpMigrations(db *sql.DB) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if err := goose.Up(db, "schema"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}

func DownMigrations(db *sql.DB) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if err := goose.Down(db, "schema"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}
