package v4_rc

import (
	"testing"
)

// TestServer_Start 测试服务器启动
func TestServer_Start(t *testing.T) {
	s := &Server{}
	err := s.Start(":8081")
	if err != nil {
		t.Fatal(err)
	}
}
