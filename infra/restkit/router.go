package restkit

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/techidea8/codectl/infra/logger"
	"github.com/techidea8/codectl/infra/utils/stringx"
)

type Router struct {
	logger                logger.ILogger
	avaiablehandlerlocker *sync.RWMutex
	avaiablehandler       map[string]map[string]reflect.Method
	avaiablemodulelocker  *sync.RWMutex
	avaiablemodule        map[string]any
}

var DefaultRouter *Router = NewRouter()

// 注册路由函数,供内部模块使用,在init 函数中进行注册
//
// ptrOfModule: 结构体指针
//
// alise...: 模型名称.如果没有填写,则以结构体名字为基础
//
// eg: register(&struct{},"module1","module2","module3")
func Register(ptrOfModule any, alise ...string) {
	DefaultRouter.Register(ptrOfModule, alise...)
}
func NewRouter() *Router {
	return &Router{
		logger:                logger.DefaultLogger,
		avaiablehandlerlocker: &sync.RWMutex{},
		avaiablehandler:       make(map[string]map[string]reflect.Method, 0),
		avaiablemodulelocker:  &sync.RWMutex{},
		avaiablemodule:        map[string]any{},
	}
}

// 注册路由
//
// alise 模型名称,如果没有填写,则以结构体名字为基础
func (r *Router) Register(ptrOfModule any, alise ...string) {
	//
	r.avaiablehandlerlocker.Lock()
	defer r.avaiablehandlerlocker.Unlock()
	ptrtype := reflect.TypeOf(ptrOfModule)
	strutype := ptrtype.Elem()
	// 没有指定,默认采用反射获取结构体名称
	if len(alise) == 0 {
		module := stringx.Ucfirst(strutype.Name())
		alise = append(alise, module)
	}
	for _, module := range alise {
		r.avaiablemodule[module] = ptrOfModule
		// 现在遍历函数
		num := ptrtype.NumMethod()
		if _, ok := r.avaiablehandler[module]; !ok {
			r.avaiablehandler[module] = make(map[string]reflect.Method, num)
		}
		for i := 0; i < num; i++ {
			met := ptrtype.Method(i)
			r.avaiablehandler[module][met.Name] = met
			if r.logger != nil {
				r.logger.Infof("regiser codeCtlHandler %s/%s=> %s\n", module, met.Name, met.Type)
			}
		}
	}
}

// @description 将 /[ctxpath]/[module]/[action] 映射到具体方法函数,
//
// @param module string 模块名称
//
// @param action string 操作名称
//
// @return ptrmodule pointer 模块名称
//
// @return method reglect.Method 方法
//
// @return err  error 错误信息
func (r *Router) Dispatch(module, action string) (ptrmodule any, method reflect.Method, err error) {
	r.avaiablehandlerlocker.Lock()
	method, ok := r.avaiablehandler[module][action]
	r.avaiablehandlerlocker.Unlock()
	if !ok {
		err = fmt.Errorf("当前服务不存在")
		return
	}
	r.avaiablemodulelocker.Lock()
	ptrmodule = r.avaiablemodule[module]
	r.avaiablemodulelocker.Unlock()
	return ptrmodule, method, nil
}