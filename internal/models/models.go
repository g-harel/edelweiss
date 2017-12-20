package models

import (
	"database/sql"
)

// Init opens a connection to the database and
// initializes the required models.
func Init() (*sql.DB, error) {
	db, err := sql.Open("postgres", `
		host=192.168.99.100
		port=5432
		user=postgres
		password=password123
		dbname=edelweiss
		sslmode=disable
	`)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		DROP TABLE IF EXISTS users;
		DROP TABLE IF EXISTS domains;

		CREATE TABLE IF NOT EXISTS domains (
			id SERIAL PRIMARY KEY,
			name VARCHAR(32) NOT NULL,
			data JSON NOT NULL,
			CONSTRAINT uq_name UNIQUE (name)
		);

		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			domain_id INTEGER NOT NULL,
			email VARCHAR(32) NOT NULL,
			hash VARCHAR(60) NOT NULL,
			CONSTRAINT uq_domain_email UNIQUE (domain_id, email),
			CONSTRAINT fk_domain_id FOREIGN KEY (domain_id)
				REFERENCES domains (id)
				ON DELETE CASCADE
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
