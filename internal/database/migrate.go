package database

import (
	"database/sql"
	"fmt"
	"strconv"

	"subscriptions-api/pkg/utils/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// ApplyMigrationsToVersion применяет миграции до указанной версии
func ApplyMigrationsToVersion(db *sql.DB, migrationsPath string, version string, log *logger.Logger) error {
	// Создаём драйвер для базы
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// Создаём мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Парсим версию
	v, err := strconv.ParseUint(version, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid version format: %w", err)
	}

	// Применяем миграции до нужной версии
	if err := m.Migrate(uint(v)); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations to version %d: %w", v, err)
	}

	log.Info("Migrations applied successfully to version", "version", v)
	return nil
}

// ApplyAllMigrations применяет все миграции
func ApplyAllMigrations(db *sql.DB, migrationsPath string, log *logger.Logger) error {
	// Создаём драйвер для базы
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// Создаём мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Применяем все миграции вверх
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Info("All migrations applied successfully")
	return nil
}
