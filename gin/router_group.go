package gee

import (
	"github.com/linyerun/mini-gin/iface"
	. "github.com/linyerun/mini-gin/utils"
	"net/http"
	"os"
)

type routerGroup struct {
	prefix      string              //分组的前缀
	middlewares []iface.HandlerFunc //支持的中间件
	parent      iface.IRouterGroup  //父分组
	engine      iface.IEngine       //所有groups共享一个engine实例
}

// 创建一个路由分组
func newRouterGroup(prefix string, parent iface.IRouterGroup, engine iface.IEngine) iface.IRouterGroup {
	return &routerGroup{
		prefix: prefix,
		parent: parent,
		engine: engine,
	}
}

func (r *routerGroup) Group(prefix string) iface.IRouterGroup {
	prefix = patternDecorate(prefix) //对前缀进行修饰，避免出问题
	newGroup := newRouterGroup(r.prefix+prefix, r, r.engine)
	r.engine.(*engine).addGroup(newGroup) //将创建的所有group加入engine的group里面
	return newGroup
}

func (r *routerGroup) GET(pattern string, handler iface.HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

func (r *routerGroup) POST(pattern string, handler iface.HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

func (r *routerGroup) Use(handler iface.IHandler) {
	r.middlewares = append(r.middlewares, newHandlerFunc(handler))
}

func (r *routerGroup) Static(relativePath, fileRoot string) { //假如是(/asserts,./static)
	relativePath = patternDecorate(relativePath)
	if len(fileRoot) > 0 && fileRoot[0] == '/' { //修饰一下fileRoot
		fileRoot = fileRoot[1:]
	}
	if _, err := os.Stat(fileRoot); err != nil { //判断文件夹是否存在
		Logger().Println("in rootGroup.Static fileRoot err:", err)
		panic(err)
	}
	handler := r.createStaticHandler(relativePath, http.Dir(fileRoot))
	r.addRoute("GET", relativePath+"/*filepath", handler) //前缀无需自己添加
}

func (r *routerGroup) addRoute(method string, pattern string, handler iface.HandlerFunc) {
	pattern = patternDecorate(pattern)                                      //对pattern进行修饰，避免出问题
	pattern = r.prefix + pattern                                            //前缀统一在这里加
	r.engine.(*engine).router.addRoute(newMethod(method), pattern, handler) //无奈之举(但是感觉很妙,框架使用者做不到)
}

func (r *routerGroup) createStaticHandler(relativePath string, fs http.FileSystem) iface.HandlerFunc {
	handler := http.StripPrefix(r.prefix+relativePath, http.FileServer(fs)) // 这个请求的前缀加上relativePath变成fs对应的前缀就是handler能处理的了
	return func(c iface.IContext) {
		if _, err := fs.Open(c.Param("filepath")); err != nil {
			Logger().Println(err)
			c.JSON(400, H{
				"code": 400,
				"msg":  "不存在该文件",
			})
			return
		}
		handler.ServeHTTP(c.GetResponse(), c.GetRequest())
	}
}
