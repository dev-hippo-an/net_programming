package tcp_data_transfer

import (
	"io"
	"net"
	"sync"
	"testing"
)

func TestProxy(t *testing.T) {
	var wg sync.WaitGroup

	// 서버 리스닝
	server, err := net.Listen("tcp", "127.0.0.1:")

	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			// 서버 accept
			conn, err := server.Accept()
			if err != nil {
				return
			}

			go func(c net.Conn) {
				defer c.Close()

				// 커넥션에 작성된 문자열을 읽어온다.
				for {
					buf := make([]byte, int32(1<<10))
					n, err := c.Read(buf)

					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}

						return
					}

					switch msg := string(buf[:n]); msg {
					case "ping":
						_, err = c.Write([]byte("pong"))
					default:
						t.Log(string(buf[:n]))
						_, err = c.Write(buf[:n])
					}

					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
				}
			}(conn)
		}
	}()

	// 클라이언트와 서버 간의 프록시 서버 셋팅 : 클라 -(1)-> 프록시 서버 -(2)-> 서버
	proxyServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			t.Log("2: Proxy before accept!")
			conn, err := proxyServer.Accept()
			t.Log("4: Proxy after accept!")
			if err != nil {
				t.Log(err)
				return
			}

			go func(from net.Conn) {
				defer from.Close()

				// 서버로 연결 (2)
				to, err := net.Dial("tcp", server.Addr().String())
				if err != nil {
					t.Error(err)
					return
				}

				defer to.Close()

				// 서버에서 클라로 클라에서 서버로 복사
				// 클라이언트와 서버가 연결되는 코드가 없이도 데이터 전송이 가능하다.
				// from - client -> proxy
				// to - proxy -> server
				err = proxy(from, to)
				if err != nil && err != io.EOF {
					t.Error(err)
				}
			}(conn)
		}
	}()

	// 클라이언트에서 프록시 서버로 연결 (1)
	t.Log("1: Client before dial!")
	conn, err := net.Dial("tcp", proxyServer.Addr().String())
	t.Log("3: Client after dial!")

	if err != nil {
		t.Fatal(err)
	}

	msgs := []struct{ Message, Replay string }{
		{"ping", "pong"},
		{"pong", "pong"},
		{"echo", "echo"},
		{"ping", "pong"},
	}

	for i, m := range msgs {
		_, err = conn.Write([]byte(m.Message))
		if err != nil {
			t.Fatal(err)
		}

		buf := make([]byte, int32(1<<10))

		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}

		actual := string(buf[:n])

		t.Logf("%q -> proxy -> %q", m.Message, actual)

		if actual != m.Replay {
			t.Errorf("%d: expected reply: %q; actual: %q", i, m.Replay, actual)
		}
	}

	t.Log("Before close all resources")
	_ = conn.Close()
	t.Log("After close connection to proxy server")
	_ = proxyServer.Close()
	_ = server.Close()
	t.Log("After close all resources")

	wg.Wait()
}
