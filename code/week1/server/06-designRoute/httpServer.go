package designRoute

import (
	"net"
	"net/http"
)

// 为确保HTTPServer结构体为Server接口的实现而定义的变量
var _ Server = &HTTPServer{}

// HTTPServer HTTP服务器
type HTTPServer struct {
	*router
}

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

// ServeHTTP WEB框架入口
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 构建上下文
	ctx := &Context{
		Req:  r,
		Resp: w,
	}

	// 查找路由树并执行命中的业务逻辑
	s.serve(ctx)

	// TODO implement me
	panic("implement me")
}

// serve 查找路由树并执行命中的业务逻辑
func (s *HTTPServer) serve(ctx *Context) {
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

// GET 注册GET请求路由
func (s *HTTPServer) GET(path string, handleFunc HandleFunc) {
	s.AddRoute(http.MethodGet, path, handleFunc)
}

// POST 注册POST请求路由
func (s *HTTPServer) POST(path string, handleFunc HandleFunc) {
	s.AddRoute(http.MethodPost, path, handleFunc)
}
