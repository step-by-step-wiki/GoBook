package staticMatchingTestCase

import "net/http"

// Server WEB服务器接口
type Server interface {
	// Handler 组合http.Handler接口
	http.Handler

	// Start 启动WEB服务器
	Start(addr string) error

	// AddRoute 注册路由
	AddRoute(method string, path string, handleFunc HandleFunc)
}
