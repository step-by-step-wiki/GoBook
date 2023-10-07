package validateParamRoute

import "net/http"

// Server WEB服务器接口
type Server interface {
	// Handler 组合http.Handler接口
	http.Handler

	// Start 启动WEB服务器
	Start(addr string) error

	// addRoute 注册路由
	addRoute(method string, path string, handleFunc HandleFunc)
}
