package main

import (
	"errors"
	"fmt"
	"mzhn/management/internal/config"

	// Библиотека для миграций

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("cannot load .env")
	}
}

func main() {
	cfg := config.New()

	pg := cfg.Pg

	cstr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", pg.User, pg.Pass, pg.Host, pg.Port, pg.Name)
	m, err := migrate.New(
		"file://migrations",
		cstr,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
	}
}
