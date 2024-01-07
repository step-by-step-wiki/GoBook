package web

import (
	"sync"
	"testing"
)

func TestContext_SafeContext(t *testing.T) {
	s := &HTTPServer{router: newRouter()}

	handleFunc := func(ctx *Context) {
		safeContext := &SafeContext{
			Context: *ctx,
			Lock:    sync.Mutex{},
		}

		type User struct {
			Name string `json:"name"`
		}

		// 获取路径参数
		id := safeContext.PathValue("name")

		user := &User{Name: id.value}

		safeContext.RespJSON(202, user)
	}

	s.GET("/user/:name", handleFunc)
	_ = s.Start(":8091")
}
