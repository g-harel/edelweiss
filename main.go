package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func initialize() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres dbname=edelweiss sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadFile("init.sql")
	if err != nil {
		panic(err)
	}

	query := string(bytes)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	return db
}

type user struct {
	ID       int    `json:"id"`
	DomainID int    `json:"domain_id"`
	Email    string `json:"email"`
	Hash     string `json:"hash"`
}

func readUsers(db *sql.DB) []user {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var res []user

	for rows.Next() {
		var (
			id       int
			domainID int
			email    string
			hash     string
		)
		rows.Scan(&id, &domainID, &email, &hash)
		res = append(res, user{id, domainID, email, hash})
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}

	return res
}

func handler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		users := readUsers(db)
		b, err := json.Marshal(users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
	}
}

func main() {
	db := initialize()
	defer db.Close()

	fmt.Printf("%+v", readUsers(db))

	router := httprouter.New()
	router.GET("/api/users", handler(db))

	http.ListenAndServe(":8080", router)
}
