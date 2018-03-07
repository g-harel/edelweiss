package session

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Store interface represents a two layer hash map.
type Store interface {
	create(id string) error
	touch(id string) error
	delete(id string) error
	get(id, key string) (string, error)
	set(id, key, value string) error
}

type store struct {
	client *redis.Client
}

// NewStore creates a new session store backed by a redis client.
func NewStore(addr, pass string) (Store, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	_, err := c.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &store{c}, nil
}

func (s *store) create(id string) error {
	p := s.client.Pipeline()

	p.HMSet(id, map[string]interface{}{"id": id})
	p.Expire(id, lifespan)

	_, err := p.Exec()
	return err
}

func (s *store) delete(id string) error {
	_, err := s.client.Del(id).Result()
	return err
}

func (s *store) touch(id string) error {
	v, err := s.client.Expire(id, lifespan).Result()
	if !v {
		return fmt.Errorf("id could not be found in store")
	}
	return err
}

func (s *store) get(id, key string) (string, error) {
	return s.client.HGet(id, key).Result()
}

func (s *store) set(id, key, value string) error {
	_, err := s.client.HSet(id, key, value).Result()
	return err
}
