package sessions

import (
	"fmt"
	"net/http"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
	"github.com/go-redis/redis"
)

// SessionID is used to unikely identify sessions.
type SessionID string

// Session contains identifying information about the user.
type Session struct {
	ip string;
	userID string;
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

// Add creates a session, saves it to the store and returns its id.
func Add(client *redis.Client, userID int, ip string) (SessionID, error) {
	hashedIP, err := bcrypt.GenerateFromPassword([]byte(ip), 6)
	if err != nil {
		return "", err
	}

	sessionID := string(rand.Int())
	status := client.HMSet(sessionID, map[string]interface{}{
		"ip": hashedIP,
		"userID": userID,
	})

	if status.Err() != nil {
		return "", status.Err()
	}

	return SessionID(sessionID), nil
}

// Get retrives a session from the session store using a session id.
func Get(client *redis.Client, sessionID SessionID) (Session, error) {
	val, err := client.HGetAll(string(sessionID)).Result()
	if err != nil {
		return Session{}, err
	}

	session := Session{
		ip: val["ip"],
		userID: val["userID"],
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