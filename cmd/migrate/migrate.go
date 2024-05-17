package main

import (
	"expenses_tracker/internal/config"
	"expenses_tracker/internal/repository"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.GetConfigFromEnv(".env")
	db, err := repository.GetSqliteDb(cfg.DB.Path)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fSrc, err := (&file.File{}).Open("./migrations")
	if err != nil {
		panic(err)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{MigrationsTable: "Migrations"})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithInstance(
		"file",
		fSrc,
		"sqlite3",
		driver,
	)
	if err != nil {
		panic(err)
	}

	m.Up()
}
