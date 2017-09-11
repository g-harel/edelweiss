package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/g-harel/edelweiss/src/models"

	"github.com/julienschmidt/httprouter"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

func connectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.99.100:6379",
		Password: "password123",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	err = client.Set("test", "true", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("test").Result()
	if err != nil {
		panic(err)
	}
	if string(val) == "true" {
		fmt.Println("✓ Redis")
	} else {
		panic(fmt.Errorf("redis test failed"))
	}

	return client
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

func authenticationMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("authenticationMiddleware")
    next.ServeHTTP(w, r)
  })
}

func main() {
	client := connectRedis()
	defer client.Close()

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

	http.ListenAndServe(":8080", authenticationMiddleware(router))
}
