package filekit

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"github.com/techidea8/codectl/infra/stringx"
)

type UploadStrategy string
type config struct {
	localdir   string
	mapperpath string
	strategy   UploadStrategy
	depth      int
	ext        string
}

var tick int64 = time.Now().Unix() % 10000

const (
	UploadStrategyDate = "date"
	UploadStrategyUUID = "uuid"
)

func defaultconfig(ext string) *config {
	return &config{
		localdir:   "/mnt/storage",
		mapperpath: "/mnt",
		depth:      2,
		strategy:   UploadStrategyUUID,
		ext:        ext,
	}
}
func (c *config) build() (filepath string, netpath string) {
	atomic.AddInt64(&tick, 1)
	filename := ""
	if c.strategy == UploadStrategyUUID {
		pk := stringx.PKID()
		arr := strings.Split(pk, "")
		arr[c.depth] = fmt.Sprintf("%s.%s", pk, c.ext)
		filename = path.Join(arr[:c.depth]...)
	} else {
		now := time.Now()
		filename = fmt.Sprintf("%d/%d/%d/%s%06d.%s", now.Year(), now.Month()+1, now.Day(), now.Format("20060102150405"), tick%10000, c.ext)
	}
	return path.Join(c.localdir, filename), path.Join(c.mapperpath, filename)
}
func SetLocalDir(dir string) UploadOption {
	return func(c *config) {
		c.localdir = dir
	}
}
func SetMapperPath(path string) UploadOption {
	return func(c *config) {
		c.mapperpath = path
	}
}

func SetDepth(dpt int) UploadOption {
	return func(c *config) {
		c.depth = dpt
	}
}
func SetStrategy(strategy UploadStrategy) UploadOption {
	return func(c *config) {
		c.strategy = strategy
	}
}
func UseStrategyDate() UploadOption {
	return func(c *config) {
		c.strategy = UploadStrategyDate
	}
}
func UseStrategyUUID() UploadOption {
	return func(c *config) {
		c.strategy = UploadStrategyUUID
	}
}

type UploadOption func(*config)

// 上传文件
func Upload(file multipart.File, header *multipart.FileHeader, opts ...UploadOption) (dstpath, filekey string, size int64, err error) {
	ext := path.Ext(header.Filename)
	c := defaultconfig(ext)
	for _, opt := range opts {
		opt(c)
	}
	dstpath, filekey = c.build()
	dir := path.Dir(dstpath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}
	dstfile, err := os.Open(dstpath)
	if err != nil {
		return
	}
	size, err = io.Copy(dstfile, file)
	return
}

// 上传文件
func UploadBase64(b64data string, ext string, opts ...UploadOption) (dstpath, filekey string, size int64, err error) {
	//文件转base64
	decodeBytes, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return
	}
	c := defaultconfig(ext)
	for _, opt := range opts {
		opt(c)
	}
	dstpath, filekey = c.build()
	dir := path.Dir(dstpath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}
	size = int64(len(decodeBytes))
	err = os.WriteFile(dstpath, decodeBytes, fs.FileMode(os.O_CREATE))
	return
}
