package tryfile

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CacheFs struct {
	cacher  Cache
	root    string
	rootdir string
}

func NewCacheFs(cacher Cache, root string) *CacheFs {
	return &CacheFs{
		cacher: cacher,
		root:   root,
	}
}
func (s *CacheFs) Walk() error {
	s.rootdir = filepath.Dir(s.root)
	return filepath.Walk(s.root, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		//_left := strings.ReplaceAll(path, pdir, "")
		//_path := filepath.Clean(filepath.ToSlash(_left))
		path = filepath.ToSlash(path)
		bts, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Println("walk", path)
		err = s.Add(path, bts, 0)
		if err != nil {
			return err
		}
		_, err = s.Get(path)
		return err
	})
}
func (fs *CacheFs) Add(key string, value []byte, tm time.Duration) error {
	return fs.cacher.Set(fs.Key(key), value, tm)
}
func (fs *CacheFs) Key(path string) string {
	return strings.ReplaceAll(filepath.ToSlash(path), filepath.ToSlash(fs.rootdir), "")
}
func (fs *CacheFs) Get(key string) (value []byte, err error) {
	return fs.cacher.Get(fs.Key(key))
}
