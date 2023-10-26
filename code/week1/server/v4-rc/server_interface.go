package v4_rc

import "net/http"

// ServerInterface 服务器实体接口
// 用于定义服务器实体的行为
type ServerInterface interface {
	// Handler 组合http.Handler接口
	http.Handler
	// Start 启动服务器
	Start(addr string) error
}
