package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/g-harel/edelweiss/src/domains"
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

func domainsHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		res, err := domains.ReadAll(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		b, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(b))
	}
}

func usersHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		res, err := users.ReadAll(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		b, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(b))
	}
}

func main() {
	db := initialize()
	defer db.Close()

	domains.Test(db)
	users.Test(db)

	router := httprouter.New()
	router.GET("/api/domains", domainsHandler(db))
	router.GET("/api/users", usersHandler(db))

	http.ListenAndServe(":8080", router)
}
