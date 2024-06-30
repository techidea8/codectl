package rediskit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCtrl struct {
	Rdb *redis.Client
}

func NewRedisCtrl(conf *RedisConf) *RedisCtrl {
	rdb := redis.NewClient(&conf.Options)
	return &RedisCtrl{
		Rdb: rdb,
	}
}

func (s *RedisCtrl) Set(k string, v interface{}, t time.Duration) error {
	cmd := s.Rdb.Set(context.Background(), k, v, t)
	return cmd.Err()
}

func (s *RedisCtrl) Get(k string) (result interface{}, err error) {
	cmd := s.Rdb.Get(context.Background(), k)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	err = cmd.Scan(&result)
	return result, err
}
