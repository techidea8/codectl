package tryfile

import (
	"bytes"
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Tryfile 实现了 http.Handler 接口，所以可以用来处理 HTTP 请求
// 其中 root 用于定义前端静态资源目录（包含js、css 文件）
// index 用于定义入口视图模板文件，通常是 index.html
// ctxpath 是应用前缀
type Tryfile struct {
	cachefs   *CacheFs
	index     string
	root      string
	ctxpath   string
	cacheable bool
	fs        *embed.FS
}

func NewTryfile() *Tryfile {
	r := &Tryfile{
		root:      "",
		index:     "index.html",
		ctxpath:   "/",
		cacheable: false,
		fs:        nil,
	}
	return r
}
func (h *Tryfile) Deploy(pkg *embed.FS) *Tryfile {
	h.fs = pkg
	return h
}

func (h *Tryfile) CtxPath(path string) *Tryfile {
	h.ctxpath = path
	return h
}

func (h *Tryfile) Root(root string) *Tryfile {
	h.root = root
	return h
}

func (h *Tryfile) Index(index string) *Tryfile {
	h.index = index
	return h
}

func (h *Tryfile) Cache() *Tryfile {
	h.cachefs = NewCacheFs(NewSyncMap(), h.root)
	go func() {
		h.cachefs.Walk()
	}()
	return h
}

func (h *Tryfile) serveEmbfsHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 路径的绝对路径
	//path  = /app1/index.html
	_path := r.RequestURI
	_path = strings.TrimPrefix(_path, h.ctxpath)
	_path = filepath.Join(h.root, _path)
	if _path == h.root {
		_path = filepath.Join(h.root, h.index)
	}
	_path = filepath.ToSlash(_path)
	bts, err := h.fs.ReadFile(_path)

	if err == nil {
		http.ServeContent(w, r, _path, time.Now(), bytes.NewReader(bts))
		return
	} else if os.IsNotExist(err) {
		bts, err := h.fs.ReadFile(filepath.ToSlash(filepath.Join(h.root, h.index)))
		if err != nil {
			// 如果期间报错，返回 500 响应
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, _path, time.Now(), bytes.NewReader(bts))
		return
	} else {
		// 如果期间报错，返回 500 响应
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 处理 Tryfile 应用请求（主要是首次访问时入口 HTML 文档和相关静态资源文件，暂不涉及 API 接口）
// https://a.b.com/app1/index.html
func (h *Tryfile) serveLocalHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 路径的绝对路径
	//path  = /app1/index.html
	_path := r.RequestURI
	_path = strings.TrimPrefix(_path, h.ctxpath)
	// 去掉前缀
	// path = /index.html
	// 在 URL 路径前加上静态资源根目录
	//path = ${root}+index.html
	_path = filepath.Join(h.root, _path)
	_, err := os.Stat(_path)
	if err == nil {
		http.StripPrefix(h.ctxpath, http.FileServer(http.Dir(h.root))).ServeHTTP(w, r)
		return
	}
	// 如果是文件不存在
	if os.IsNotExist(err) {
		_index := filepath.ToSlash(filepath.Join(h.root, h.index))
		http.ServeFile(w, r, _index)
		return
	} else {
		// 如果期间报错，返回 500 响应
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Tryfile) serveLocalCacheHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 路径的绝对路径
	//path  = /app1/index.html
	_path := r.RequestURI
	_path = strings.TrimPrefix(_path, h.ctxpath)
	// 去掉前缀
	// path = /index.html
	// 在 URL 路径前加上静态资源根目录
	//path = ${root}+index.html
	_path = filepath.Join(h.root, _path)
	_path = filepath.ToSlash(_path)
	// 检查对应资源文件是否存在
	//D:/winlion/zkxdr_integrated_machine_base/assist/console/js/index.7ucGj7Xq.js
	bts, err := h.cachefs.Get(_path)
	if err == nil {
		http.ServeContent(w, r, _path, time.Now(), bytes.NewReader(bts))
		return
	} else if os.IsNotExist(err) {
		_index := filepath.Join(h.root, h.index)
		bts, err := h.cachefs.Get(_index)
		if err != nil {
			// 如果期间报错，返回 500 响应
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, _path, time.Now(), bytes.NewReader(bts))
		return
	} else {
		// 如果期间报错，返回 500 响应
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (h *Tryfile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.fs != nil {
		h.serveEmbfsHTTP(w, r)
	} else {
		if h.cachefs != nil {
			h.serveLocalCacheHTTP(w, r)
		} else {
			h.serveLocalHTTP(w, r)
		}
	}
}
