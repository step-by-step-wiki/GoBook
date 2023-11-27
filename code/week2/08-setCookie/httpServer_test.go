package setCookie

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	// Tips: 初始化HTTPServer结构体时,不能将其类型声明为Server接口的实现,因为一些方法被我们定义在了
	// Tips: 实现上,而不是在接口上,所以不能将其类型声明为Server接口的实现.即:不能写成如下形式:
	// var s Server = &HTTPServer{}
	s := &HTTPServer{router: newRouter()}

	// 注册1个处理函数到路由
	handleFoo := func(ctx *Context) { fmt.Println("处理第1件事") }
	handleBar := func(ctx *Context) { fmt.Println("处理第2件事") }

	// 将2个处理函数封装为1个处理函数
	handleAssemble := func(ctx *Context) {
		handleFoo(ctx)
		handleBar(ctx)
	}

	s.addRoute(http.MethodGet, "/getUser", handleAssemble)

	_ = s.Start(":8080")
}

func TestServer_serve(t *testing.T) {
	s := NewHTTPServer()
	handleFunc := func(ctx *Context) {
		// 直接调用http.ResponseWriter的Write方法时 默认响应码为200
		ctx.Resp.Write([]byte("hello order detail"))
	}
	s.GET("/order/detail", handleFunc)

	wildcardHandleFunc := func(ctx *Context) {
		respPath := ""

		pathSegments := strings.Split(ctx.Req.URL.Path, "/")
		for _, pathSegment := range pathSegments {
			if pathSegment == "" {
				continue
			}
			respPath += pathSegment
			respPath += " "
		}

		respPath = strings.TrimRight(respPath, " ")
		respMsg := fmt.Sprintf("hello %s", respPath)

		ctx.Resp.Write([]byte(respMsg))
	}
	s.GET("/order/*", wildcardHandleFunc)

	_ = s.Start(":8081")
}
