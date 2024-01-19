package recover_panic

import (
	"fmt"
	"testing"
	"web"
)

// Test_MiddlewareBuilder 测试捕获panic中间件构造器
func Test_MiddlewareBuilder(t *testing.T) {
	// 创建捕获panic中间件构造器
	builder := &MiddlewareBuilder{
		StatusCode: 500,
		Data:       []byte("panic error"),
		LogFunc: func(ctx *web.Context) {
			fmt.Printf("panic路径: %s\n", ctx.Req.URL.Path)
		},
	}

	// 构建捕获panic中间件
	option := web.ServerWithMiddleware(builder.Build())

	// 创建web服务器
	server := web.NewHTTPServer(option)

	// 创建HandleFunc
	handleFunc := func(ctx *web.Context) {
		panic("test panic")
	}

	// 注册路由并启动服务器
	server.GET("/user", handleFunc)
	server.Start(":8080")
}
