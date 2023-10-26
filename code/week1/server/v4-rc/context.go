package v4_rc

import "net/http"

// Context 路由处理函数的上下文
type Context struct {
	Req        *http.Request       // Req HTTP请求
	Resp       http.ResponseWriter // Resp HTTP响应
	PathParams map[string]string   // PathParams 参数路由的参数
}
