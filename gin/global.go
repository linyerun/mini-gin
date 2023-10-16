package gee

import (
	"html/template"
	"io"
	"reflect"
)

var (
	//前端渲染(不想因为这个就要把engine整合到每个context里面)
	//直接放在全局, 让它在engine中被配置, 在context的HTML中被使用
	funcMap      template.FuncMap
	htmlTemplate *template.Template

	// UserLogWriter 用户自定义的io输出流
	userLogWriter io.Writer
)

// SetMiniGinLogOutputWriter 允许用户自定义日志的输出流
func SetMiniGinLogOutputWriter(logWriter io.Writer) {
	// 判断一下这个值, 避免为nil值
	isNil := reflect.ValueOf(logWriter).IsNil()
	if isNil {
		panic("can not set logWriter to be nil")
	}
	// 设置
	userLogWriter = logWriter
}

// GetMiniGinLogOutputWriter 获取用户自定义日志的输出流
func GetMiniGinLogOutputWriter() io.Writer {
	return userLogWriter
}
