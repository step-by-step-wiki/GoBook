package web

import (
	"net"
	"net/http"
)

// 为确保HTTPServer结构体为Server接口的实现而定义的变量
var _ Server = &HTTPServer{}

// HTTPServer HTTP服务器
type HTTPServer struct {
	router                   // router 路由树
	middlewares []Middleware // middlewares 中间件切片.表示HTTPServer需要按顺序执行的的中间件链
}

// NewHTTPServer 创建HTTP服务器
// 这里选项的含义其实是指不同的 Option 函数
// 每一个 Option 函数都会对 HTTPServer 实例的不同属性进行设置
func NewHTTPServer(opts ...Option) *HTTPServer {
	server := &HTTPServer{
		router: newRouter(),
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

// NewHTTPServerWithConfig 从给定的配置文件中读取配置项 根据配置项创建HTTP服务器
func NewHTTPServerWithConfig(configFilePath string) *HTTPServer {
	// 在函数体内读取配置文件 解析为一个配置对象
	// 然后根据配置对象创建 HTTPServer 实例
	return nil
}

// ServerWithMiddleware 本函数用于为 HTTPServer 实例添加中间件
// 即: 本函数用于设置 HTTPServer 实例的middlewares属性
func ServerWithMiddleware(middlewares ...Middleware) Option {
	return func(server *HTTPServer) {
		server.middlewares = middlewares
	}
}

// ServeHTTP WEB框架入口
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 构建上下文
	ctx := &Context{
		Req:  r,
		Resp: w,
	}

	// 执行中间件链
	root := s.serve
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		root = s.middlewares[i](root)
	}

	// 查找路由树并执行命中的业务逻辑
	root(ctx)
}

// serve 查找路由树并执行命中的业务逻辑
func (s *HTTPServer) serve(ctx *Context) {
	method := ctx.Req.Method
	path := ctx.Req.URL.Path
	targetNode, ok := s.findRoute(method, path)
	// 没有在路由树中找到对应的路由节点 或 找到了路由节点的处理函数为空(即NPE:none pointer exception 的问题)
	// 则返回404
	if !ok || targetNode.node.HandleFunc == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		// 此处确实会报错 但是作为一个WEB框架 遇上了这种错误也没有特别好的处理办法
		// 最多只能是落个日志
		_, _ = ctx.Resp.Write([]byte("Not Found"))
		return
	}

	// 命中节点则将路径参数名值对设置到上下文中
	ctx.PathParams = targetNode.pathParams

	// 命中节点则将节点的路由设置到上下文中
	ctx.MatchRoute = targetNode.node.route

	// 执行路由节点的处理函数
	targetNode.node.HandleFunc(ctx)
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
	s.addRoute(http.MethodGet, path, handleFunc)
}

// POST 注册POST请求路由
func (s *HTTPServer) POST(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodPost, path, handleFunc)
}
