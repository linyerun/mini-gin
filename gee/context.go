package gee

import (
	"encoding/json"
	"fmt"
	"github.com/linyerun/Gee/iface"
	. "github.com/linyerun/Gee/utils"
	"net/http"
)

type context struct {
	// origin objects
	resp http.ResponseWriter
	req  *http.Request

	// request info
	path   string
	method string
	params map[string]string

	// response info
	statusCode int

	//实现中间件
	index    int                 //当前执行的中间件(-1表示在0的前面)
	handlers []iface.HandlerFunc //当前的handler放在这个的最后
}

func newContext(r *http.Request, w http.ResponseWriter) iface.IContext {
	return &context{
		req:    r,
		resp:   w,
		path:   r.URL.Path,
		method: r.Method,
		index:  -1,
	}
}

// 获取更细度的接口

func (c *context) GetRequest() *http.Request {
	return c.req
}

func (c *context) GetResponse() http.ResponseWriter {
	return c.resp
}

func (c *context) GetPath() string {
	return c.path
}

func (c *context) GetMethod() string {
	return c.method
}

func (c *context) GetStatusCode() int {
	return c.statusCode
}

// 获取参数接口(还应该加一个获取Path的和Body的才行)

func (c *context) PostForm(key string) string {
	//TODO 对这个方法掌握程度不够
	return c.req.FormValue(key)
}

func (c *context) Query(key string) string {
	return c.req.URL.Query().Get(key)
}

func (c *context) Param(key string) string {
	return c.params[key]
}

//设置返回状态

func (c *context) Status(code int) {
	c.statusCode = code
	c.resp.WriteHeader(code)
}

//设置响应行

func (c *context) SetHeader(key string, value string) {
	c.resp.Header().Set(key, value)
}

//设置返回数据格式和内容

func (c *context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)

	if _, err := c.resp.Write([]byte(fmt.Sprintf(format, values...))); err != nil {
		panic(err)
	}
}

func (c *context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.resp)
	if err := encoder.Encode(obj); err != nil {
		Logger().Errorln("Gee=> encoder.Encode err:", err)
		http.Error(c.resp, err.Error(), 500)
	}
}

func (c *context) Data(code int, data []byte) {
	c.Status(code)
	if _, err := c.resp.Write(data); err != nil {
		panic(err)
	}
}

func (c *context) HTML(code int, name string, data any) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := htmlTemplate.ExecuteTemplate(c.resp, name, data); err != nil {
		Logger().Errorln("Your HTML has error:", err)
		c.String(500, "system has err,code =", 500)
		panic(err)
	}
}

//实现中间件
func (c *context) next() {
	c.index++
	c.handlers[c.index](c) //作者这里使用的是循环，因为有的用户不写next,或者说没用next的逻辑，但是我提供的接口自己执行next,保证了无需循环。
}
