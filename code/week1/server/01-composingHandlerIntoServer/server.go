package composingHandlerIntoServer

import "net/http"

type Server interface {
	// Handler 组合http.Handler接口
	http.Handler
}
