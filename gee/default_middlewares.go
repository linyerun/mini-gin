package gee

import (
	"fmt"
	"github.com/linyerun/Gee/iface"
	. "github.com/linyerun/Gee/utils"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type loggerHandler struct {
	t time.Time
}

func (f *loggerHandler) PrevHandle(_ iface.IContext) {
	f.t = time.Now()
}

func (f *loggerHandler) LastHandle(c iface.IContext) {
	Logger().Printf("%s [%d] %s in %v", c.GetRequest().Host, c.GetStatusCode(), c.GetRequest().RequestURI, time.Since(f.t))
}

type recoverHandler struct {
}

func (r *recoverHandler) PrevHandle(c iface.IContext) {
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			Logger().Printf("%s\n\n", trace(message))
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}()
}

func (r *recoverHandler) LastHandle(_ iface.IContext) {
}

// print stack trace for debug (调试用的)
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
