package cachekit

import (
	"context"
	"time"

	"github.com/techidea8/codectl/infra/rediskit"
)

type rediscache struct {
	ctrl *rediskit.RedisCtrl
}

func NewRedisCache(conf *rediskit.RedisConf) *rediscache {
	ctrl := rediskit.NewRedisCtrl(conf)
	return &rediscache{
		ctrl: ctrl,
	}
}
func (r *rediscache) Set(k string, v interface{}, d time.Duration) error {
	cmd := r.ctrl.Rdb.Set(context.Background(), k, v, d)
	return cmd.Err()
}

// 获得
func (r *rediscache) Get(k string) (result interface{}, err error) {
	cmd := r.ctrl.Rdb.Get(context.Background(), k)
	if err = cmd.Err(); err != nil {
		return nil, err
	}
	err = cmd.Scan(result)
	return result, err
}
