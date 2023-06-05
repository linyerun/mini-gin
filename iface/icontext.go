package iface

import "net/http"

//为操作req和resp做封装

type HandlerFunc func(IContext)

type IContext interface {
	GetRequest() *http.Request
	GetResponse() http.ResponseWriter
	GetPath() string
	GetMethod() string
	GetStatusCode() int
	PostForm(key string) string
	Query(key string) string
	Param(key string) string
	Status(code int)
	SetHeader(key string, value string)
	String(code int, format string, values ...interface{})
	JSON(code int, obj interface{})
	Data(code int, data []byte)
	HTML(code int, name string, data any)
}
