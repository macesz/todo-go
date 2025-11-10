package postgres

import (
	"embed"
	"errors"
	"fmt"
	"log"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

func MigrateDb(DbUser, DbPass, dbAddr, DbName string) error {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)

	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		return err
	}

	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("migration done")

	return nil
}
