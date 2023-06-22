package gee

import (
	"github.com/linyerun/Gee/iface"
	. "github.com/linyerun/Gee/utils"
	"html/template"
	"net/http"
	"strings"
)

type engine struct {
	iface.IRouterGroup
	router *router
	groups []iface.IRouterGroup
}

func New() iface.IEngine {
	e := &engine{router: newRouter()}
	e.IRouterGroup = newRouterGroup("", nil, e)
	e.addGroup(e.IRouterGroup) //e本身也是一个IRouter, 所以也加入groups里面
	e.Use(new(recoverHandler))
	return e
}

func Default() iface.IEngine {
	e := New().(*engine)
	e.Use(new(loggerHandler)) // 看一个打进了的请求运行了多少秒
	return e
}

// SetFuncMap 前端相关
func (e *engine) SetFuncMap(f template.FuncMap) {
	funcMap = f
}

func (e *engine) LoadHTMLGlob(templatePathPattern string) {
	Logger().Infoln("You should be sure use SetFuncMap before use LoadHTMLGlob!")
	glob, err := template.New("global").Funcs(funcMap).ParseGlob(templatePathPattern)
	if err != nil {
		panic(err)
	}
	htmlTemplate = glob
}

func (e *engine) Run(addr string) (err error) {
	Logger().Debugln("Gee Running.....")
	return http.ListenAndServe(addr, e)
}

//所有请求的入口方法
func (e *engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	Logger().Debugln(req.Host, "coming....", "request:[", req.Method, "]", "[", req.URL.Path, "]")
	c := newContext(req, resp).(*context)
	//开始将满足要求的middleware加入到c里面
	for _, group := range e.groups {
		g := group.(*routerGroup)
		if strings.HasPrefix(c.GetPath(), g.prefix) {
			c.handlers = append(c.handlers, g.middlewares...)
		}
	}
	e.router.handle(c)
}

func (e *engine) addGroup(group iface.IRouterGroup) {
	e.groups = append(e.groups, group)
}
