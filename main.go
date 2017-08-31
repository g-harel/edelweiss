package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/g-harel/edelweiss/src/users"

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

func handler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		res := users.ReadAll(db)
		b, err := json.Marshal(res)
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

	// adding users
	err := users.Add(db, &users.User{
		DomainID: 1,
		Email: "email1@example.com",
		Hash: "password123",
	})
	if err != nil {
		fmt.Println(err)
	}
	err = users.Add(db, &users.User{
		DomainID: 1,
		Email: "email2@example.com",
		Hash: "password123",
	})
	if err != nil {
		fmt.Println(err)
	}
	err = users.Add(db, &users.User{
		DomainID: 2,
		Email: "email1@example.com",
		Hash: "password123",
	})
	if err != nil {
		fmt.Println(err)
	}
	//

	router := httprouter.New()
	router.GET("/api/users", handler(db))

	http.ListenAndServe(":8080", router)
}
