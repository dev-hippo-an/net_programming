package tcp_data_transfer

import (
	"crypto/rand"
	"io"
	"net"
	"testing"
)

func TestReadIntoBuffer(t *testing.T) {

	// 16MB 의 랜덤 데이터 생성
	payload := make([]byte, 1<<24)
	_, err := rand.Read(payload)

	if err != nil {
		t.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}

		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				t.Error(err)
				return
			}
		}(conn)

		_, err = conn.Write(payload)
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())

	if err != nil {
		t.Fatal(err)
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			t.Error(err)
			return
		}
	}(conn)

	// 클라이언트(tcp 에선 클라이언트 개념이 없이 노드이지만 dial 로 연결한 쪽을 클라이언트로 한다.) 버퍼 사이즈 약 512 KB
	buf := make([]byte, 1<<19)

	for {
		// 반복적으로 16MB 의 페이로드를 클라이언트 버퍼 사이즈 크기로 분할되어 읽는다.
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}

		t.Logf("Read : %d byte", n)
	}

}
