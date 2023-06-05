package gee

import (
	"github.com/linyerun/Gee/iface"
	"strings"
)

func patternDecorate(pattern string) string {
	//对pattern修饰一下
	pattern = strings.TrimSuffix(pattern, "/")
	if len(pattern) == 0 || pattern[0] != '/' {
		pattern = "/" + pattern
	}
	return pattern
}

func getParts(pattern string) (parts []string) {
	for _, part := range strings.Split(pattern, "/") {
		if len(part) != 0 {
			parts = append(parts, part)
		}
	}
	return
}

func parsePattern(pattern string, parts []string) map[string]string {
	m := make(map[string]string)
	p := getParts(pattern)
	for i, s := range p {
		if strings.HasPrefix(s, ":") {
			m[s[1:]] = parts[i]
		} else if strings.HasPrefix(s, "*") {
			m[s[1:]] = "/" + strings.Join(parts[i:], "/")
		}
	}
	return m
}

func newHandlerFunc(handler iface.IHandler) iface.HandlerFunc {
	return func(c iface.IContext) {
		//模板方法模式
		handler.PrevHandle(c)
		c.(*context).next()
		handler.LastHandle(c)
	}
}
