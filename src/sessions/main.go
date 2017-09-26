package sessions

import (
	"fmt"
	"net/http"
	"math/rand"
	"strconv"

	"github.com/go-redis/redis"
)

type SessionID string

type Session struct {
	ip string;
	userID int;
}

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

func Add(client *redis.Client, userID int, ip string) (SessionID, error) {
	// TODO hash ip
	sessionID := string(rand.Int())
	status := client.HMSet(sessionID, map[string]interface{}{
		"ip": ip,
		"userID": userID,
	})

	if status.Err() != nil {
		return "", status.Err()
	}

	return SessionID(sessionID), nil
}

func Get(client *redis.Client, sessionID SessionID) (Session, error) {
	// TODO hardcoded
	session := Session{}
	val, err := client.HGet(string(sessionID), "ip").Result()
	if err != nil {
		return session, err
	}

	session.ip = val

	val, err = client.HGet(string(sessionID), "userID").Result()
	if err != nil {
		return session, err
	}

	session.userID, err = strconv.Atoi(val)
	if err != nil {
		return session, err
	}

	return session, nil
}

func Middleware() func(http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("authenticationMiddleware")
			next.ServeHTTP(w, r)
		})
	}
}