package database

import (
	"database/sql"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

// New opens a connection to the database and runs migrations.
func New(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// for development. re-creates db every restart
	migrations := &migrate.FileMigrationSource{
		Dir: "services/gateway/database/migrations",
	}
	migrate.SetTable("migrations")
	if true {
		m, err := migrate.Exec(db, "postgres", migrations, migrate.Down)
		if err == nil {
			fmt.Printf("Downed %d migrations!\n", m)
		}
	}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return nil, err
	}
	if n > 0 {
		fmt.Printf("Applied %d migrations!\n", n)
	}

	return db, nil
}
