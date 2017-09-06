package models

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

// IDomains is the domains' database interface.
type IDomains interface {
	Add(name string, data string) (int, error)
	UpdateData(id int, data string) error
	ReadAll() ([]Domain, error)
}

// Domains database
type Domains struct {
	DB *sql.DB
}

// Add adds a new domain to the database and returns the id.
func (d Domains) Add(name string, data string) (int, error) {
	stmt, err := d.DB.Prepare(`
		INSERT INTO domains (name, data)
			VALUES ($1, $2)
			RETURNING (id)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = validateName(name)
	if err != nil {
		return 0, err
	}

	var id int
	row := stmt.QueryRow(name, data)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateData will change the data for a domain in the database.
// Authentication must be done before this point.
func (d Domains) UpdateData(id int, data string) error {
	stmt, err := d.DB.Prepare(`
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
func (d Domains) ReadAll() ([]Domain, error) {
	rows, err := d.DB.Query(`
		SELECT * FROM domains
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := readDomains(rows)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func readDomains(rows *sql.Rows) ([]Domain, error) {
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
