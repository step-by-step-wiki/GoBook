package addRoutes

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	s := &HTTPServer{}

	// 注册多个处理函数到路由
	handleFoo := func(ctx Context) { fmt.Println("处理第1件事") }
	handleBar := func(ctx Context) { fmt.Println("处理第2件事") }

	s.AddRoutes(http.MethodGet, "/getUser", handleFoo, handleBar)

	_ = s.Start(":8080")
}
