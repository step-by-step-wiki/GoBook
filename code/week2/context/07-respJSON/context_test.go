package respJSON

import (
	"testing"
)

func TestContext_RespJSON(t *testing.T) {
	s := &HTTPServer{router: newRouter()}

	handleFunc := func(ctx *Context) {
		type User struct {
			Name string `json:"name"`
		}

		// 获取路径参数
		id := ctx.PathValue("name")

		user := &User{Name: id.value}

		ctx.RespJSON(202, user)
	}

	s.GET("/user/:name", handleFunc)
	_ = s.Start(":8091")
}
