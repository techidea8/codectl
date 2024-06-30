package cachekit

import "time"

//接口
type Cache interface {
	// 设置
	Set(k string, x interface{}, d time.Duration) error
	// 获得
	Get(k string) (interface{}, error)
}

var DefaultCache Cache = NewMemoryCache(time.Minute * 10)

func UseCache(cache Cache) {
	DefaultCache = cache
}

func Set(k string, v interface{}, d time.Duration) error {
	return DefaultCache.Set(k, v, d)
}

func Get(k string) (interface{}, error) {
	return DefaultCache.Get(k)
}
