package tcp_session

import (
	"io"
	"net"
	"testing"
)

// 리스너와 핸들러를 정의, Dial 을 통한 요청
func TestDial(t *testing.T) {
	// 로컬 호스트의 랜덤 포트에 tcp 연결 리스너 바인딩
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer func() {
			done <- struct{}{}
		}()

		for {
			conn, err := listener.Accept()
			t.Log("connection accepted")

			if err != nil {
				t.Log(err)
				return
			}

			// handler
			go func(c net.Conn) {
				defer func() {
					c.Close() // tcp 세션 graceful terminate
					done <- struct{}{}
				}()

				// 소켓으로부터 1024바이트를 읽어 수신 데이터 로깅
				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()

	// 리스너에 연결 시도
	conn, err := net.Dial("tcp", listener.Addr().String())

	if err != nil {
		t.Fatal(err)
	}

	conn.Close() // tcp 세션 graceful terminate
	<-done
	listener.Close()
	<-done
}
