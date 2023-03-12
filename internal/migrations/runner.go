package migrations

import (
	"database/sql"
	"fmt"
	"shopping-cart-backend/config"
	"shopping-cart-backend/pkg/database"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

// Up is used for running migrations on DB
func Up(dbConfig config.Database) error {
	databaseCfg := database.Config{
		Name:     config.App.Database.Name,
		User:     config.App.Database.User,
		Password: config.App.Database.Password,
		Host:     config.App.Database.Host,
		Port:     config.App.Database.Port,
		SSL:      config.App.Database.SSLMode,
	}
	db, err := sql.Open(dbConfig.Dialect, databaseCfg.ConnectionURL())

	if err != nil {
		return fmt.Errorf("connection to Postgres failed: %s", err.Error())
	}

	migrations := &migrate.FileMigrationSource{
		Dir: dbConfig.MigrationsDir,
	}

	_, err = migrate.Exec(db, dbConfig.Dialect, migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("migrations up failed due to: %v", err)
	}
	return nil
}

// Down is used for destroying migrations in a DB, if any
func Down(dbConfig config.Database) error {
	databaseCfg := database.Config{
		Name:     config.App.Database.Name,
		User:     config.App.Database.User,
		Password: config.App.Database.Password,
		Host:     config.App.Database.Host,
		Port:     config.App.Database.Port,
		SSL:      config.App.Database.SSLMode,
	}
	db, err := sql.Open(dbConfig.Dialect, databaseCfg.ConnectionURL())

	if err != nil {
		return fmt.Errorf("connection to Postgres failed: %s", err.Error())
	}

	migrations := &migrate.FileMigrationSource{
		Dir: dbConfig.MigrationsDir,
	}

	_, err = migrate.Exec(db, dbConfig.Dialect, migrations, migrate.Down)
	if err != nil {
		return fmt.Errorf("migrations down failed due to: %v", err)
	}
	return nil
}
