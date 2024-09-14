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
type UploadConf struct {
	LocalDir   string
	MapperPath string
	Strategy   UploadStrategy
	Depth      int
	ext        string
}

var tick int64 = time.Now().Unix() % 10000

const (
	UploadStrategyDate = "date"
	UploadStrategyUUID = "uuid"
)

func defaultconfig(ext string) *UploadConf {
	return &UploadConf{
		LocalDir:   "/mnt/storage",
		MapperPath: "/mnt",
		Depth:      2,
		Strategy:   UploadStrategyUUID,
		ext:        ext,
	}
}
func (c *UploadConf) build() (filepath string, netpath string) {
	atomic.AddInt64(&tick, 1)
	filename := ""
	if c.Strategy == UploadStrategyUUID {
		pk := stringx.PKID()
		arr := strings.Split(pk, "")
		arr[c.Depth] = fmt.Sprintf("%s%s", pk, c.ext)
		filename = path.Join(arr[:c.Depth+1]...)
	} else {
		now := time.Now()
		filename = fmt.Sprintf("%d/%d/%d/%s%06d%s", now.Year(), now.Month()+1, now.Day(), now.Format("20060102150405"), tick%10000, c.ext)
	}
	return path.Join(c.LocalDir, filename), path.Join(c.MapperPath, filename)
}
func SetLocalDir(dir string) UploadOption {
	return func(c *UploadConf) {
		c.LocalDir = dir
	}
}
func SetMapperPath(path string) UploadOption {
	return func(c *UploadConf) {
		c.MapperPath = path
	}
}

func SetDepth(dpt int) UploadOption {
	return func(c *UploadConf) {
		c.Depth = dpt
	}
}
func SetStrategy(strategy UploadStrategy) UploadOption {
	return func(c *UploadConf) {
		c.Strategy = strategy
	}
}
func UseStrategyDate() UploadOption {
	return func(c *UploadConf) {
		c.Strategy = UploadStrategyDate
	}
}
func UseStrategyUUID() UploadOption {
	return func(c *UploadConf) {
		c.Strategy = UploadStrategyUUID
	}
}

type UploadOption func(*UploadConf)

// 上传文件
func Upload(file multipart.File, header *multipart.FileHeader, opts ...UploadOption) (dstpath, filekey, filename, ext string, size int64, err error) {
	filename = header.Filename
	ext = path.Ext(header.Filename)
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
	dstfile, err := os.Create(dstpath)
	if err != nil {
		return
	}
	defer dstfile.Close()
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
