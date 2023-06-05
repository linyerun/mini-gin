package gee

import "html/template"

var (
	//前端渲染(不想因为这个就要把engine整合到每个context里面)
	//直接放在全局, 让它在engine中被配置, 在context的HTML中被使用
	funcMap      template.FuncMap
	htmlTemplate *template.Template
)
