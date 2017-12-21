package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/g-harel/edelweiss/internal/models"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

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
	db, err := models.Init()
	defer db.Close()
	if err != nil {
		panic(err)
	}

	router := httprouter.New()

	users := models.Users{DB: db}
	domains := models.Domains{DB: db}

	router.GET("/api/users", sendJSON(users.ReadAll()))

	router.GET("/api/domains", sendJSON(domains.ReadAll()))

	http.ListenAndServe(":8080", sessions.Middleware()(router))
}
