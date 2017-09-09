package models

import (
	"fmt"
	"database/sql"
)

var db *sql.DB

var Domains IDomains
var Users IUsers

func connect() error {
	psql, err := sql.Open("postgres", `
		host=192.168.99.100
		port=5432
		user=postgres
		password=password123
		dbname=edelweiss
		sslmode=disable
	`)
	if err != nil {
		return err
	}

	err = psql.Ping()
	if err != nil {
		return err
	}

	db = psql
	return nil
}

func Init() error {
	err := connect()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		DROP TABLE IF EXISTS users;
		DROP TABLE IF EXISTS domains;
	`)
	if err != nil {
		return err
	}

	domains, err := CreateDomains(db)
	if err != nil {
		return err
	}

	users, err := CreateUsers(db)
	if err != nil {
		return err
	}

	err = TestDomains(domains)
	if err != nil {
		fmt.Println(err)
	}
	err = TestUsers(users)
	if err != nil {
		fmt.Println(err)
	}

	Domains = domains
	Users = users

	return nil
}

func Close() error {
	return db.Close()
}