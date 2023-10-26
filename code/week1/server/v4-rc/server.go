package v4_rc

import (
	"net"
	"net/http"
)

var _ ServerInterface = &HTTPServer{}

type HTTPServer struct {
	*router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

// ServeHTTP 是http.Handler接口的方法 此处必须先写个实现 不然Server不是http.Handler接口的实现
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Req:  r,
		Resp: w,
	}

	s.serve(ctx)
}

// serve 查找路由树并执行匹配到的节点所对应的处理函数
func (s *HTTPServer) serve(ctx *Context) {
	method := ctx.Req.Method
	path := ctx.Req.URL.Path
	targetNode, found := s.router.findRoute(method, path)
	if !found || targetNode.node.HandleFunc == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("not found"))
		return
	}
	ctx.PathParams = targetNode.pathParams
	targetNode.node.HandleFunc(ctx)
}

// Start 启动服务器
func (s *HTTPServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return http.Serve(listener, s)
}

// GET 注册GET路由
func (s *HTTPServer) GET(path string, handler HandleFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

// POST 注册POST路由
func (s *HTTPServer) POST(path string, handler HandleFunc) {
	s.addRoute(http.MethodPost, path, handler)
}
