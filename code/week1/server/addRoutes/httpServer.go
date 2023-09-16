package addRoutes

import (
	"net"
	"net/http"
)

// 为确保HTTPServer结构体为Server接口的实现而定义的变量
var _ Server = &HTTPServer{}

// HTTPServer HTTP服务器
type HTTPServer struct {
}

// ServeHTTP WEB框架入口
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO implement me
	panic("implement me")
}

// Start 启动WEB服务器
func (s *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// 在监听端口之后,启动服务之前做一些操作
	// 例如在微服务框架中,启动服务之前需要注册服务

	return http.Serve(l, s)
}

// AddRoute 注册路由
func (s *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {
	// TODO: implement me
	panic("implement me")
}

// AddRoutes 支持1个路由对应多个处理函数的注册路由
func (s *HTTPServer) AddRoutes(method string, path string, handles ...HandleFunc) {
	// TODO: implement me
	panic("implement me")
}
