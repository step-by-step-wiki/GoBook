package v6

import (
	"fmt"
	"net/http"
	"testing"
)

// Test_Middleware 测试中间件的工作顺序
func Test_Middleware(t *testing.T) {
	s := NewHTTPServer()

	s.middlewares = []Middleware{
		Middleware1,
		Middleware2,
		Middleware3,
		Middleware4,
	}

	s.ServeHTTP(nil, &http.Request{})
}

func Middleware1(next HandleFunc) HandleFunc {
	return func(ctx *Context) {
		fmt.Println("中间件1开始执行")
		next(ctx)
		fmt.Println("中间件1结束执行")
	}
}

func Middleware2(next HandleFunc) HandleFunc {
	return func(ctx *Context) {
		fmt.Println("中间件2开始执行")
		next(ctx)
		fmt.Println("中间件2结束执行")
	}
}

func Middleware3(next HandleFunc) HandleFunc {
	return func(ctx *Context) {
		fmt.Println("中间件3中断后续中间件的执行")
	}
}

func Middleware4(next HandleFunc) HandleFunc {
	return func(ctx *Context) {
		fmt.Println("中间件4不会被执行")
	}
}
