package access_log

import (
	"fmt"
	"net/http"
	"testing"
	"web"
)

// Test_MiddlewareBuilder 以创建 http.Request 的方式,测试记录日志中间件
func Test_MiddlewareBuilder(t *testing.T) {
	// 创建记录日志中间件
	builder := &MiddlewareBuilder{}
	// 注意这里的链式调用 使得在创建中间件的同时就可以设置记录日志的函数
	accessLogMiddleware := builder.SetLogFunc(func(logContent string) {
		fmt.Println(logContent)
	}).Build()

	// 创建中间件Option
	middlewareOption := web.ServerWithMiddleware(accessLogMiddleware)

	// 创建服务器
	// 这里就可以看出Option模式的好处了: 通过不同的Option函数可以设置HTTPServer实例的不同属性
	s := web.NewHTTPServer(middlewareOption)

	// 注册路由
	s.GET("/a/b/*", func(ctx *web.Context) {
		fmt.Println("hello, it's me")
	})

	// 创建请求
	request, err := http.NewRequest(http.MethodGet, "/a/b/c", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 启动服务器
	s.ServeHTTP(nil, request)
}

// Test_MiddlewareBuilderWithServer 以启动服务器的方式测试记录日志中间件
func Test_MiddlewareBuilderWithServer(t *testing.T) {
	// 创建记录日志中间件
	builder := &MiddlewareBuilder{}
	// 注意这里的链式调用 使得在创建中间件的同时就可以设置记录日志的函数
	accessLogMiddleware := builder.SetLogFunc(func(logContent string) {
		fmt.Println(logContent)
	}).Build()

	// 创建中间件Option
	middlewareOption := web.ServerWithMiddleware(accessLogMiddleware)

	// 创建服务器
	// 这里就可以看出Option模式的好处了: 通过不同的Option函数可以设置HTTPServer实例的不同属性
	s := web.NewHTTPServer(middlewareOption)

	// 注册路由
	s.GET("/a/b/*", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, it's me"))
	})

	// 启动服务器
	s.Start(":8092")
}
