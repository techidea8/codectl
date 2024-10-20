package filekit

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"github.com/techidea8/codectl/infra/stringx"
)

type StorageStrategy string
type StorageConf struct {
	LocalDir   string
	MapperPath string
	Strategy   StorageStrategy
	ServerUrl  string //服务地址
	Depth      int
	ext        string
}

var tick int64 = time.Now().Unix() % 10000

const (
	StorageStrategyDate = "date"
	StorageStrategyUUID = "uuid"
)

func defaultconfig(ext string) *StorageConf {
	return &StorageConf{
		LocalDir:   "/mnt/storage",
		MapperPath: "/mnt",
		Depth:      2,
		ServerUrl:  "",
		Strategy:   StorageStrategyUUID,
		ext:        ext,
	}
}
func (c *StorageConf) build() (filepath string, netpath string) {
	atomic.AddInt64(&tick, 1)
	filename := ""
	if c.Strategy == StorageStrategyUUID {
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
func SetLocalDir(dir string) StorageOption {
	return func(c *StorageConf) {
		c.LocalDir = dir
	}
}
func SetMapperPath(path string) StorageOption {
	return func(c *StorageConf) {
		c.MapperPath = path
	}
}

func SetServerUrl(url string) StorageOption {
	return func(c *StorageConf) {
		c.ServerUrl = url
	}
}
func SetDepth(dpt int) StorageOption {
	return func(c *StorageConf) {
		c.Depth = dpt
	}
}
func SetStrategy(strategy StorageStrategy) StorageOption {
	return func(c *StorageConf) {
		c.Strategy = strategy
	}
}
func UseStrategyDate() StorageOption {
	return func(c *StorageConf) {
		c.Strategy = StorageStrategyDate
	}
}
func UseStrategyUUID() StorageOption {
	return func(c *StorageConf) {
		c.Strategy = StorageStrategyUUID
	}
}

type StorageOption func(*StorageConf)

// 上传文件
func Realpath(filekey string, opts ...StorageOption) (realpath string) {
	conf := defaultconfig("")
	for _, opt := range opts {
		opt(conf)
	}
	localdir := conf.LocalDir
	mappath := conf.MapperPath
	realpath = strings.ReplaceAll(filekey, mappath, localdir)
	return
}

// 上传文件
func ServeFile(filekey string, opts ...StorageOption) http.Handler {
	realpath := Realpath(filekey, opts...)
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.ServeFile(rw, req, realpath)
	})
}

// 上传文件
func Upload(file multipart.File, header *multipart.FileHeader, opts ...StorageOption) (dstpath, filekey, filename, ext string, size int64, err error) {
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
func StorageBase64(b64data string, ext string, opts ...StorageOption) (dstpath, filekey string, size int64, err error) {
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
