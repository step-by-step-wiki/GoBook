package v4_rc

import (
	"net"
	"net/http"
)

type Server struct {
}

// ServeHTTP 是http.Handler接口的方法 此处必须先写个实现 不然Server不是http.Handler接口的实现
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

// Start 启动服务器
func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return http.Serve(listener, s)
}

// addRoute 注册路由
func (s *Server) addRoute(method string, pattern string, handler HandleFunc) {
	panic("implement me")
}

// GET 注册GET路由
func (s *Server) GET(path string, handler HandleFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

// POST 注册POST路由
func (s *Server) POST(path string, handler HandleFunc) {
	s.addRoute(http.MethodPost, path, handler)
}
