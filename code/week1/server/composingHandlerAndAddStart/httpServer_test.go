package composingHandlerAndAddStart

import "testing"

func TestServer(t *testing.T) {
	s := &HTTPServer{}
	s.Start(":8084")
}
