package domains

import (
	"database/sql"
	"errors"
	"regexp"
)

// Domain is a domain
type Domain struct {
	ID   int    `json:"-"`
	Name string `json:"email"`
	Data string `json:"data"`
}

func validateName(name string) error {
	matched, err := regexp.Match("^\\w{2,32}$", []byte(name))
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("Invalid name")
	}

	return nil
}

func read(rows *sql.Rows) ([]Domain, error) {
	res := []Domain{}
	for rows.Next() {
		domain := new(Domain)
		err := rows.Scan(&domain.ID, &domain.Name, &domain.Data)
		if err != nil {
			continue
		}
		res = append(res, *domain)
	}

	err := rows.Err()
	if err != nil {
		return res, err
	}

	return res, nil
}

// Add adds a new domain to the database.
// The new domain's id is added to the Domain struct.
func Add(db *sql.DB, d *Domain) error {
	stmt, err := db.Prepare(`
		INSERT INTO domains (name, data)
			VALUES ($1, $2)
			RETURNING (id)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = validateName(d.Name)
	if err != nil {
		return err
	}

	row := stmt.QueryRow(d.Name, d.Data)
	err = row.Scan(&d.ID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateData will change the data for a domain in the database.
// Authentication is must be done before this point.
func UpdateData(db *sql.DB, id int, data string) error {
	stmt, err := db.Prepare(`
		UPDATE domains SET data=$1 WHERE id=$2
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(data, id)
	if err != nil {
		return err
	}

	return nil
}

// ReadAll reads all domains from the database
func ReadAll(db *sql.DB) ([]Domain, error) {
	rows, err := db.Query(`
		SELECT * FROM domains
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := read(rows)
	if err != nil {
		return nil, err
	}

	return res, nil
}
