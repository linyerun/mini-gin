package iface

import "html/template"

//添加路由、启动程序入口

type IEngine interface {
	IRouterGroup
	Run(addr string) (err error)
	SetFuncMap(f template.FuncMap)
	LoadHTMLGlob(templatePathPattern string)
}
