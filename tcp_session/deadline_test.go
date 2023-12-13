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
			conn.Close()
			close(sync) // 이른 return 으로 인해 sync 채널에서 읽는 데이터가 블로킹되면 안 됨
		}()

		err = conn.SetDeadline(time.Now().Add(5 * time.Second)) // 1
		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1)
		_, err = conn.Read(buf) // 원격 노드가 데이터를 보낼 때까지 블로킹

		nErr, ok := err.(net.Error)
		if !ok || !nErr.Timeout() { // 2
			t.Errorf("expected timeout error; actual: %v", err)
		}
		sync <- struct{}{}

		err = conn.SetDeadline(time.Now().Add(5 * time.Second)) // 3
		if err != nil {
			t.Error(err)
			return
		}

		_, err = conn.Read(buf)
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()
	<-sync
	_, err = conn.Write([]byte("1"))

	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err != io.EOF { // 4
		t.Errorf("expected server termination; actual: %v", err)
	}

}