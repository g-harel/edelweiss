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
}