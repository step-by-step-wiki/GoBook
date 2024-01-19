package access_log

// accessLog 定义访问日志的结构
type accessLog struct {
	Host       string `json:"host,omitempty"`       // Host 主机地址
	Route      string `json:"route,omitempty"`      // Route 命中的路由
	HTTPMethod string `json:"HTTPMethod,omitempty"` // HTTPMethod 请求的HTTP方法
	Path       string `json:"path,omitempty"`       // Path 请求的路径 即请求的uri部分
}
