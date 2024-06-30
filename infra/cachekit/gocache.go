package cachekit

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type gocache struct {
	cacher *cache.Cache
}

func NewMemoryCache(defaultduration time.Duration) *gocache {
	cacher := cache.New(defaultduration, 2*defaultduration)
	return &gocache{
		cacher: cacher,
	}
}
func (s *gocache) Set(k string, v interface{}, d time.Duration) error {
	s.cacher.Set(k, v, d)
	return nil
}

// 获得
func (s *gocache) Get(k string) (interface{}, error) {
	r, f := s.cacher.Get(k)
	if !f {
		return nil, fmt.Errorf("获取失败")
	} else {
		return r, nil
	}
}
