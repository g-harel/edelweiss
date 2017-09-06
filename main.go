package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/g-harel/edelweiss/src/models"

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

type generator func() (interface{}, error)

func sendJSON(payload generator) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		res, err := payload()
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

	domains := models.Domains{DB: db}
	users := models.Users{DB: db}

	err := models.TestDomains(domains)
	if err != nil {
		fmt.Println(err)
	}
	err = models.TestUsers(users)
	if err != nil {
		fmt.Println(err)
	}

	router := httprouter.New()

	router.GET("/api/domains", sendJSON(func() (interface{}, error) {
		return domains.ReadAll()
	}))

	router.GET("/api/users", sendJSON(func() (interface{}, error) {
		return users.ReadAll()
	}))

	http.ListenAndServe(":8080", router)
}
