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

func (s *store) create(ID string) error {
	pipe := s.client.Pipeline()

	pipe.HMSet(ID, map[string]interface{}{"ID": ID})
	pipe.Expire(ID, lifespan)

	_, err := pipe.Exec()
	return err
}

func (s *store) delete(ID string) error {
	_, err := s.client.Del(ID).Result()
	return err
}

func (s *store) touch(ID string) error {
	_, err := s.client.Expire(ID, lifespan).Result()
	return err
}

func (s *store) get(ID, key string) (string, error) {
	return s.client.HGet(ID, key).Result()
}

func (s *store) set(ID, key, value string) error {
	_, err := s.client.HSet(ID, key, value).Result()
	return err
}
