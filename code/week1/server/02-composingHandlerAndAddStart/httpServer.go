package composingHandlerAndAddStart

import (
	"net"
	"net/http"
)

type HTTPServer struct {
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO implement me
	panic("implement me")
}

func (s *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// 在监听端口之后,启动服务之前做一些操作
	// 例如在微服务框架中,启动服务之前需要注册服务

	return http.Serve(l, s)
}
