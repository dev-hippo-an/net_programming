package tcp_session

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestDeadline(t *testing.T) {
	sync := make(chan struct{})

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

		defer func() {
			t.Log("Let's close connection in go routine")
			conn.Close()
			close(sync) // 이른 return 으로 인해 sync 채널에서 읽는 데이터가 블로킹되면 안 됨
		}()

		err = conn.SetDeadline(time.Now().Add(7 * time.Second)) // 1 - 연결된 connection 에 Deadline 설정
		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1)
		_, err = conn.Read(buf) // 원격 노드가 데이터를 보낼 때까지 블로킹 -> 타임 아웃이 발생할겁니다.

		nErr, ok := err.(net.Error)
		if !ok || !nErr.Timeout() { // 2 - DL 로 인한 timeout ERR
			t.Errorf("expected timeout error; actual: %v", err)
		}

		sync <- struct{}{} // 채널에 데이터를 보내고 난 후

		t.Log("Timeout~@@@@ Let's get back to code flow in go routine")

		err = conn.SetDeadline(time.Now().Add(1 * time.Second)) // 3
		if err != nil {
			t.Error(err)
			return
		}

		_, err = conn.Read(buf) // 데이터 수신 블로킹
		if err != nil {
			t.Errorf("here is second timeout err: %s", err)
		}
		t.Log("Read from connection")
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()
	<-sync // go routine 에서 채널에 보낸 데이터를 받을 때까지 수신 대기 (buffer capacity 0)

	t.Log("let's write to connection")
	_, err = conn.Write([]byte("1")) // 채널 데이터를 수신 후 블로킹이 해제되어 코드 진행

	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1)
	_, err = conn.Read(buf) // connection 에서 원격 노드가 데이터를 보낼 때까지 블로킹 -> connection 이 닫히고 EOF 반환

	t.Logf("Connection closed here; actual: %v", err)
	if err != io.EOF { // 4
		t.Errorf("expected server termination; actual: %v", err)
	}

}
