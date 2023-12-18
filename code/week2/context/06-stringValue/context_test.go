package stringValue

import (
	"fmt"
	"testing"
)

func TestContext_PathValue(t *testing.T) {
	s := &HTTPServer{router: newRouter()}

	handleFunc := func(ctx *Context) {
		// 获取路径参数
		id, err := ctx.PathValue("id").AsInt64()
		if err != nil {
			ctx.Resp.WriteHeader(400)
			ctx.Resp.Write([]byte("id输入不正确: " + err.Error()))
			return
		}

		ctx.Resp.Write([]byte(fmt.Sprintf("hello %d", id)))
	}

	s.GET("/order/:id", handleFunc)
	_ = s.Start(":8091")
}
