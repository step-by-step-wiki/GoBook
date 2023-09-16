package httpServer

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	http.ListenAndServe(":8085", nil)
}
