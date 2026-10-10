package postgres

import (
	"errors"
	"fmt"
	"strings"
	"todo-service/internal/storage/postgres/migrations"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func newMigrator(databaseUrl string) (*migrate.Migrate, error) {
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("migrations source: %w", err)
	}

	url := strings.Replace(databaseUrl, "postgres://", "pgx5://", 1)

	m, err := migrate.NewWithSourceInstance("iofs", src, url)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return m, nil
}

func RunMigrations(databaseUrl string, command string) error {
	m, err := newMigrator(databaseUrl)
	if err != nil {
		return fmt.Errorf("new migrator: %w", err)
	}

	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			return
		}
	}(m)

	switch command {
	case "up":
		err = m.Up()
	case "down":
		err = m.Steps(-1)
	case "version":
		version, dirty, versionErr := m.Version()
		if versionErr != nil {
			return versionErr
		}

		fmt.Printf("version: %d, dirty: %v\n", version, dirty)
		return nil
	default:
		return fmt.Errorf("unknown command %q (use up, down, version)", command)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
