package restkit

import (
	"net/http"
	"regexp"
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
type RegexCfg struct {
	patern  *regexp.Regexp
	param   []string
	handerx *HandlerFuncX
}
type Router struct {
	parent         *Router
	prefix         string
	logger         logger.ILogger
	pathMap        *sync.Map
	regMap         *sync.Map
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
		regMap:         &sync.Map{},
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
		regMap:         &sync.Map{},
		handleNotFound: http.NotFound,
	}
}
func (h *Router) PathPrefix(prefix string) *Router {
	prefix = strings.TrimPrefix(prefix, "/")
	prefix = "/" + prefix
	h.prefix = prefix
	return h
}

func (h *Router) NotFound(handler http.Handler) {
	h.handleNotFound = handler.ServeHTTP
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
	patern := strings.Join(slicekit.Reverse(tmpprefix), "")
	// 如果 包含 {},则认为是规则匹配
	// /a/b/c/{d}/e

	if strings.Contains(patern, "{") {
		params, regpatern := h.buildparamandreg(patern)
		_patern := regexp.MustCompile(regpatern)
		cfg := &RegexCfg{
			patern:  _patern,
			param:   params,
			handerx: hander,
		}
		rootrouter.regMap.Store(patern, cfg)
		//_reg
	} else {
		rootrouter.pathMap.Store(patern, hander)
	}
	return hander
}

// 根据 路径生成参数表和正则字符串
func (h *Router) buildparamandreg(patern string) (params []string, regpatern string) {
	keys := []string{}
	index := 0
	reg := []rune{}
	lock := false
	for i, c := range patern {
		if c == '{' {
			index = i + 1
			lock = true
			reg = append(reg, '(', '.', '*', ')')
		} else if c == '}' {
			keys = append(keys, patern[index:i])
			lock = false
		} else {
			if !lock {
				reg = append(reg, c)
			}
		}
	}
	return keys, string(reg)
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
	uri := req.RequestURI
	_input := strings.Split(uri, "?")[0]
	hf, ok := h.pathMap.Load(_input)
	ctx := NewContext(w, req)
	if !ok {
		// 如果没找到,那么正则去匹配
		h.regMap.Range(func(key, value any) bool {
			cfg := value.(*RegexCfg)
			arr := cfg.patern.FindStringSubmatch(_input)
			// 找到
			if len(arr) > 1 {
				hf = cfg.handerx
				// 找到
				lenarr := len(arr)
				for i := 0; i < len(cfg.param); i++ {
					paramName := cfg.param[i]
					if i < lenarr {
						paramValue := arr[i+1]
						ctx.AddParam(paramName, paramValue)
					} else {
						ctx.AddParam(paramName, "")
					}
				}
				return false
			}
			return true
		})
		if hf == nil {
			h.handleNotFound(w, req)
			return
		}
	}

	handlerfuncx := hf.(*HandlerFuncX)
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
		// for i := len(h.postmiddleware) - 1; i >= 0; i-- {
		// 	_handlerfuncx = h.postmiddleware[i](_handlerfuncx)
		// }
		// _handlerfuncx(ctx)
	}
}
