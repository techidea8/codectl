package restkit

import (
	"net/http"
	"strings"
	"sync"

	"github.com/techidea8/codectl/infra/logger"
	"github.com/techidea8/codectl/infra/slicekit"
	"github.com/techidea8/codectl/infra/wraper"
)

type HandlerFuncX struct {
	handler HandlerFunc
	methods []string
}
type Router struct {
	parent         *Router
	prefix         string
	logger         logger.ILogger
	pathMap        *sync.Map
	handleNotFound http.HandlerFunc
	premiddleware  []MiddlewareFunc
	postmiddleware []MiddlewareFunc
}

func NewRouter() *Router {
	return &Router{
		parent:         nil,
		prefix:         "/",
		logger:         logger.DefaultLogger,
		pathMap:        &sync.Map{},
		handleNotFound: http.NotFound,
		premiddleware:  make([]MiddlewareFunc, 0),
		postmiddleware: make([]MiddlewareFunc, 0),
	}
}
func (h *Router) Subrouter() *Router {

	return &Router{
		parent:         h,
		prefix:         "/",
		logger:         logger.DefaultLogger,
		pathMap:        &sync.Map{},
		handleNotFound: http.NotFound,
	}
}
func (h *Router) PathPrefix(prefix string) *Router {
	prefix = strings.TrimPrefix(prefix, "/")
	prefix = "/" + prefix
	h.prefix = prefix
	return h
}
func (h *Router) HandleFunc(path string, fun HandlerFunc) *HandlerFuncX {
	path = strings.TrimPrefix(path, "/")
	path = "/" + path
	hander := &HandlerFuncX{
		handler: fun,
		methods: []string{
			http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodDelete, http.MethodPut,
		},
	}
	tmpprefix := []string{
		path,
	}
	rootrouter := h
	for h.parent != nil {
		if h.prefix != "" {
			tmpprefix = append(tmpprefix, h.prefix)
		}
		h = h.parent
		rootrouter = h
	}
	if rootrouter.prefix != "" && rootrouter.prefix != "/" {
		tmpprefix = append(tmpprefix, rootrouter.prefix)
	}
	rootrouter.pathMap.Store(strings.Join(slicekit.Reverse(tmpprefix), ""), hander)
	return hander
}

// 使用日志
func (h *Router) UseLogger(logger logger.ILogger) *Router {
	h.logger = logger
	return h
}

// 使用日志
func (h *Router) Pre(middle ...MiddlewareFunc) *Router {
	h.premiddleware = append(h.premiddleware, middle...)
	return h
}

// 使用日志
func (h *Router) Post(middle ...MiddlewareFunc) *Router {
	h.postmiddleware = append(h.postmiddleware, middle...)
	return h
}

// 支持method
func (h *HandlerFuncX) Methods(method ...string) {
	tmps := []string{}
	for _, v := range method {
		tmps = append(tmps, strings.ToUpper(v))
	}
	h.methods = append([]string{}, tmps...)
}

var DefaultRouter *Router = NewRouter()

// 提供服务
func (h *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	hf, ok := h.pathMap.Load(req.RequestURI)
	if !ok {
		h.handleNotFound(w, req)
	} else {
		handlerfuncx := hf.(*HandlerFuncX)
		ctx := NewContext(w, req)
		// 如果不包含
		if !slicekit.Contains(handlerfuncx.methods, req.Method) {
			h.handleNotFound(w, req)
		} else {
			_handlerfuncx := handlerfuncx.handler
			for i := len(h.premiddleware) - 1; i >= 0; i-- {
				_handlerfuncx = h.premiddleware[i](_handlerfuncx)
			}

			result, err := _handlerfuncx(ctx)
			if err != nil {
				wraper.Error(err).Encode(w)
			} else {
				result.Encode(w)
			}
			for i := len(h.postmiddleware) - 1; i >= 0; i-- {
				_handlerfuncx = h.postmiddleware[i](_handlerfuncx)
			}
			_handlerfuncx(ctx)
		}
	}
}
