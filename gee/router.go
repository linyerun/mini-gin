package gee

import (
	"fmt"
	"github.com/linyerun/Gee/iface"
	. "github.com/linyerun/Gee/utils"
	"net/http"
)

type method string

func newMethod(m string) method {
	return method(m)
}

func (m method) getName() string {
	return string(m)
}

type router struct {
	handlers map[string]iface.HandlerFunc //pattern=>handler
	roots    map[method]*node             //动态路由(GET、POST...各对应一个node)
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]iface.HandlerFunc),
		roots:    make(map[method]*node),
	}
}

func (r *router) addRoute(method method, pattern string, handler iface.HandlerFunc) {
	pattern = patternDecorate(pattern) //对pattern修饰一下

	key := method.getName() + "_" + pattern
	r.handlers[key] = handler

	root, ok := r.roots[method]
	if !ok {
		root = new(node)
		r.roots[method] = root
	}

	parts := getParts(pattern) //获取各个部分

	root.insert(pattern, parts, 0)

	Logger().Debugf("Route %4s - %s", method, pattern)
}

func (r *router) handle(c iface.IContext) {
	//对pattern进行修饰
	pattern := patternDecorate(c.GetPath())

	//通过method获取对应的root
	methodName := c.GetMethod()
	root, ok := r.roots[newMethod(methodName)]
	if !ok {
		//找不到
		c.JSON(http.StatusNotFound, map[string]any{
			"code": http.StatusNotFound,
			"msg":  fmt.Sprintf("404 NOT FOUND,path: %s、method: %s!", pattern, methodName),
		})
		return
	}
	parts := getParts(pattern) //得到parts

	node := root.search(parts, 0) //进行搜索
	if node == nil {
		//找不到
		c.JSON(http.StatusNotFound, map[string]any{
			"code": http.StatusNotFound,
			"msg":  fmt.Sprintf("404 NOT FOUND,path: %s、method: %s!", pattern, methodName),
		})
		return
	}

	realCtx := c.(*context)

	realCtx.params = parsePattern(node.pattern, parts) //这样更妙(用户改不了params值)

	handlerFunc := r.handlers[methodName+"_"+node.pattern]
	realCtx.handlers = append(realCtx.handlers, handlerFunc)

	realCtx.next() //开始执行
}
