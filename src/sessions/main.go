package sessions

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
)

// Init opens a connection to the client.
func Init() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.99.100:6379",
		Password: "password123",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	// TODO
	Test(client)

	return client, nil
}

func Middleware() func(http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("authenticationMiddleware")
			next.ServeHTTP(w, r)
		})
	}
}