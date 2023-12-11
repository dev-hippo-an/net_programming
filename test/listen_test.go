package test

import (
	"net"
	"testing"
)

// net.Listen 함수를 수신 연결 요청 처리가 가능한 TCP 서버(리스너) 작성
func TestListener(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = listener.Close()
	}()

	t.Logf("bound to %q", listener.Addr())
}
