package sessions

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Test runs some basic tests on the client.
func Test(client *redis.Client) {
	err := client.Set("test", "true", 0).Err()
	if err != nil {
		fmt.Println("redis test failed")
	}

	val, err := client.Get("test").Result()
	if err != nil {
		fmt.Println("redis test failed")
	}
	if string(val) == "true" {
		fmt.Println("✓ Redis")
	} else {
		fmt.Println("redis test failed")
	}

	id, err := Add(client, 1234567890, "192.0.0.1")
	if err != nil {
		fmt.Println(err)
	}

	session, err := Get(client, id)
	if err != nil {
		fmt.Println(err)
	}

	if session.userID == "1234567890" && session.ip != "192.0.0.1" {
		fmt.Println("✓ Sessions")
	}
}