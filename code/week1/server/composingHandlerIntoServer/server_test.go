package composingHandlerIntoServer

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	var s Server
	http.ListenAndServe(":8085", s)
	http.ListenAndServeTLS(":443", "", "", s)
}
