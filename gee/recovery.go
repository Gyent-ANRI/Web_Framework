package gee

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

//以中间件的形式加入错误恢复
func Recovery() HandleFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				weberr.Printf("Panic : %v\n\tTrace: %v\n", err, trace())
				ctx.HTML(http.StatusInternalServerError, "", "Internal server error")
			}
		}()
		ctx.Next()
	}
}

//追溯发生panic的位置
func trace() string {
	var pcs [32]uintptr
	var str strings.Builder
	//跳过callers, trace和defer三个调用以简化Log
	n := runtime.Callers(3, pcs[:])
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%v: %v", file, line))
	}
	return str.String()
}
