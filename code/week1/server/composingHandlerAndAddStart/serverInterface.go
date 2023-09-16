package composingHandlerAndAddStart

import "net/http"

type Server interface {
	http.Handler
	Start(addr string) error
}
