package session

import (
	"github.com/go-redis/redis"
)

type store struct {
	client *redis.Client
}

func newStore() (*store, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.99.100:6379",
		Password: "password123",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &store{client}, nil
}

func (s *store) create(id string) error {
	pipe := s.client.Pipeline()

	pipe.HMSet(id, map[string]interface{}{"id": id})
	pipe.Expire(id, lifespan)

	_, err := pipe.Exec()
	return err
}

func (s *store) delete(id string) error {
	_, err := s.client.Del(id).Result()
	return err
}

func (s *store) touch(id string) error {
	_, err := s.client.Expire(id, lifespan).Result()
	return err
}

func (s *store) get(id, key string) (string, error) {
	return s.client.HGet(id, key).Result()
}

func (s *store) set(id, key, value string) error {
	_, err := s.client.HSet(id, key, value).Result()
	return err
}

func (s *store) close() error {
	return s.client.Close()
}
