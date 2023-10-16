package main

import (
	"fmt"
	"github.com/linyerun/mini-gin/gin"
	"github.com/linyerun/mini-gin/iface"
	"html/template"
	"time"
)

func main() {
	e := gee.Default()

	// 普通路由方法的使用
	e.GET("/lin", func(c iface.IContext) {
		c.JSON(200, gee.H{"msg": "GET method,ok!", "code": 200})
	})
	e.POST("/ye", func(c iface.IContext) {
		c.JSON(200, gee.H{"msg": "POST method,ok!", "code": 200})
	})

	// 加载HTML文件的位置, 配合c.HTML使用
	e.SetFuncMap(template.FuncMap{"FormatAsDate": FormatAsDate}) //设置参数要在加载模板前进行
	e.LoadHTMLGlob("templates/*")                                // 加载所有的html文件

	e.GET("/html", func(c iface.IContext) {
		c.HTML(200, "format_as_date.html", gee.H{
			"title": "gin",
			"now":   time.Date(2023, 06, 22, 16, 20, 0, 0, time.UTC),
		})
	})

	// 加载静态文件
	e.Static("/assets", "./static")

	// 使用中间件
	e.Use(new(MyMiddleware))

	// 使用路由分组功能
	userGroup := e.Group("/user")
	{
		userGroup.GET("/lin", func(c iface.IContext) {
			c.JSON(200, gee.H{"msg": "GET method,ok!", "code": 200})
		})
		userGroup.POST("/ye", func(c iface.IContext) {
			c.JSON(200, gee.H{"msg": "POST method,ok!", "code": 200})
		})
		userGroup.Static("/assets", "./user_file")
		e.Use(new(MyMiddleware))
	}

	// 运行
	panic(e.Run(":8080"))
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

type MyMiddleware struct {
}

func (m *MyMiddleware) PrevHandle(c iface.IContext) {
	fmt.Println(c.GetPath()+"执行前时间:", time.Now().Format("2006-01-02 15-04-05"))
}

func (m *MyMiddleware) LastHandle(c iface.IContext) {
	fmt.Println(c.GetPath()+"执行后时间:", time.Now().Format("2006-01-02 15-04-05"))
}
