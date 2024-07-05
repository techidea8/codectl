package restkit

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	"github.com/techidea8/codectl/infra/logger"
	"github.com/techidea8/codectl/infra/utils/stringx"
	"github.com/techidea8/codectl/infra/wraper"
)

type pathHandler struct {
	ctxpath      string
	logger       logger.ILogger
	router       *Router
	beforemiddle []MiddlewareFunc
	aftermiddle  []MiddlewareFunc
}

func NewHandler(ctxpath string) *pathHandler {
	return &pathHandler{
		ctxpath:      ctxpath,
		router:       DefaultRouter,
		beforemiddle: make([]MiddlewareFunc, 0),
		aftermiddle:  make([]MiddlewareFunc, 0),
	}
}

// 使用日志
func (h *pathHandler) UseLogger(logger logger.ILogger) *pathHandler {
	h.logger = logger
	return h
}

// 使用日志
func (h *pathHandler) Pre(middle ...MiddlewareFunc) *pathHandler {
	h.beforemiddle = append(h.beforemiddle, middle...)
	return h
}

// 使用日志
func (h *pathHandler) Post(middle ...MiddlewareFunc) *pathHandler {
	h.aftermiddle = append(h.aftermiddle, middle...)
	return h
}

// 配置router
func (h *pathHandler) UseRouter(router *Router) *pathHandler {
	h.router = router
	return h
}

// 网络服务
func (h pathHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 通过反射获得,远不如原生高效
	ctx := &context{
		request: req,
		writer:  w,
	}
	handle := func(req Context) (r *wraper.Response, err error) {
		//CTXPATH/:module/:action
		paternstr := h.ctxpath + `/([^/^?]+)/([^/^?]+)\??([^?]*)`
		if h.ctxpath == "" || h.ctxpath == "/" {
			paternstr = `/([^/^?]+)/([^/^?]+)\??([^?]*)`
		}
		patern := regexp.MustCompile(paternstr)
		result := patern.FindAllStringSubmatch(req.Request().RequestURI, -1)
		if len(result) < 1 {
			return wraper.Error("当前服务" + req.Request().RequestURI + "不存在").HttpStatus(http.StatusNotFound), fmt.Errorf("当前服务" + req.Request().RequestURI + "不存在")
		}
		arr := result[0]
		// 找不到系统服务
		if len(arr) < 3 {
			return wraper.Error("当前服务" + req.Request().RequestURI + "不存在").HttpStatus(http.StatusNotFound), fmt.Errorf("当前服务" + req.Request().RequestURI + "不存在")
		}
		module := stringx.Ucfirst(arr[1])
		action := stringx.Ucfirst(arr[2])
		ptrmodule, method, e := h.router.Dispatch(module, action)
		if e != nil {
			return wraper.Error(e).HttpStatus(http.StatusNotFound), fmt.Errorf("当前服务" + e.Error())
		}
		if ptrmodule == nil {
			return wraper.Error(e).HttpStatus(http.StatusNotFound), fmt.Errorf("当前服务" + req.Request().RequestURI + "不存在")
		}
		// diyige canshu  shi jiegouti
		reply := method.Func.Call([]reflect.Value{reflect.ValueOf(ptrmodule), reflect.ValueOf(req)})
		resultwraper := reply[0].Interface()
		err = nil
		if reply[1].Interface() != nil {
			err = reply[1].Interface().(error)
		}
		if err != nil {
			return wraper.Error(err), err
		} else {
			if resultwraper == nil {
				logger.Infof("empty when %s", req.Request().RequestURI)
				return wraper.Empty(), nil
			} else {
				return resultwraper.(*wraper.Response), nil
			}
		}
	}

	for i := len(h.beforemiddle) - 1; i >= 0; i-- {
		handle = h.beforemiddle[i](handle)
	}
	result, err := handle(ctx)
	if err != nil {
		wraper.Error(err).Encode(w)
	} else {
		result.Encode(w)
	}
	for i := len(h.aftermiddle) - 1; i >= 0; i-- {
		handle = h.aftermiddle[i](handle)
	}

}
