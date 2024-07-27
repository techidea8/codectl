package tryfile

import (
	"os"
	"sync"
	"time"

	"github.com/coocood/freecache"
)

type Cache interface {
	Set(key string, value []byte, tm time.Duration) error
	Get(key string) ([]byte, error)
}

type freecacher struct {
	cache *freecache.Cache
}

func NewFreecacher() *freecacher {
	//30M
	return &freecacher{
		cache: freecache.NewCache(10 << 20),
	}
}

type syncmap struct {
	datamap sync.Map
}

func NewSyncMap() *syncmap {
	return &syncmap{}
}
func (f *syncmap) Set(key string, value []byte, tm time.Duration) error {
	f.datamap.Store(key, value)
	return nil
}
func (f *syncmap) Get(key string) (value []byte, err error) {
	ret, ok := f.datamap.Load(key)
	if ok {
		return ret.([]byte), nil
	} else {
		return nil, os.ErrNotExist
	}
}

func (f *freecacher) Set(key string, value []byte, tm time.Duration) error {
	f.cache.Set([]byte(key), value, -1)
	return nil
}
func (f *freecacher) Get(key string) (value []byte, err error) {
	return f.cache.Get([]byte(key))
}
