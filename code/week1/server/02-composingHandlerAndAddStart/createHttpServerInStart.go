package composingHandlerAndAddStart

import (
	"net"
	"net/http"
)

type HTTPServerWithOriginHttpServer struct {
}

func (s *HTTPServerWithOriginHttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO implement me
	panic("implement me")
}

func (s *HTTPServerWithOriginHttpServer) Start(addr string) error {
	originServer := http.Server{
		Addr:    addr,
		Handler: s,
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return originServer.Serve(l)
}
