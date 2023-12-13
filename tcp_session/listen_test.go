package tcp_session

import (
	"net"
	"testing"
)

// net.Listen 함수를 수신 연결 요청 처리가 가능한 TCP 서버(리스너) 작성
func TestListener(t *testing.T) {
	// listener 가 특정 ip 주소와 포트에 바인딩된다.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = listener.Close()  // 리스너를 우아하게 종료하기 위해 사용
	}()

	t.Logf("bound to %q", listener.Addr())
}
