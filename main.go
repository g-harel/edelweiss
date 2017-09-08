package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/g-harel/edelweiss/src/models"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func connectPostgres() *sql.DB {
	db, err := sql.Open("postgres", `
		host=192.168.99.100
		port=5432
		user=postgres
		password=password123
		dbname=edelweiss
		sslmode=disable
	`)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func sendJSON(res interface{}, err error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
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
	db := connectPostgres()
	defer db.Close()

	_, err := db.Exec(`
		DROP TABLE IF EXISTS users;
		DROP TABLE IF EXISTS domains;
	`)
	if err != nil {
		panic(err)
	}

	domains, err := models.CreateDomains(db)
	if err != nil {
		panic(err)
	}

	users, err := models.CreateUsers(db)
	if err != nil {
		panic(err)
	}

	err = models.TestDomains(domains)
	if err != nil {
		fmt.Println(err)
	}
	err = models.TestUsers(users)
	if err != nil {
		fmt.Println(err)
	}

	router := httprouter.New()

	router.GET("/api/domains", sendJSON(domains.ReadAll()))

	router.GET("/api/users", sendJSON(users.ReadAll()))

	http.ListenAndServe(":8080", router)
}
