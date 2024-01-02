package web

import "net/http"

// Server WEB服务器接口
type Server interface {
	http.Handler                                                // Handler 组合http.Handler接口
	Start(addr string) error                                    // Start 启动WEB服务器
	addRoute(method string, path string, handleFunc HandleFunc) // addRoute 注册路由
}
