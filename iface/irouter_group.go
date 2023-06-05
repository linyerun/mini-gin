package iface

type IRouterGroup interface {
	Group(prefix string) IRouterGroup         //创建分组
	GET(pattern string, handler HandlerFunc)  //分组的GET方法
	POST(pattern string, handler HandlerFunc) //分组的POST方法
	Use(handler IHandler)
	Static(relativePath, fileRoot string)
}
