package tcp_session

import (
	"context"
	"io"
	"net"
	"testing"
	"time"
)

func TestPingAdvanceDeadline(t *testing.T) {
	done := make(chan struct{})

	// 리스너 설정
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	begin := time.Now()
	go func() {
		defer func() {
			close(done)
		}()

		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())

		defer func() {
			cancel()
			conn.Close()
		}()

		resetTimer := make(chan time.Duration, 1) // cap 1개 블로킹 없이 코드 진행
		resetTimer <- time.Second
		go Pinger(ctx, conn, resetTimer)

		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1024)

		for {
			// connection 에 write 되는 타이밍은 처음 4번의 read 가 종료되고 난 후 - 90번째 라인 코드
			// 4초가 된 후에 write 가 발생하고 read 가 되면서 아래에서 connection 의 데드라인을 다시 5초로 설정
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			t.Logf("In listener accept go routine -> [%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])

			//resetTimer <- 0

			// 5초의 타임아웃이 설정되고나서 done 채널이 close 됨
			err = conn.SetDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				t.Error(err)
				return
			}
		}
	}()

	// 리스너에 연결 설정
	conn, err := net.Dial("tcp", listener.Addr().String())

	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	buf := make([]byte, 1024)

	for w := 0; w < 3; w++ {
		for i := 0; i < 4; i++ {
			// ping.go 의 타이머 만료로 ping 이 connection 에 write 된다.\
			// 해당 값을 읽어 오는 것
			n, err := conn.Read(buf)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("In first ping test - [%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
		}

		_, err = conn.Write([]byte("Pong!"))
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 4; i++ {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
			break
		}
		t.Logf("In second ping test -> [%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])

	}

	<-done // 읽어올 값이 없어 대기하다가 done 채널이 close 가 된 후 (위의 고루틴의 자원 정리 함수에서 진행) 흐름 진행
	end := time.Since(begin).Truncate(time.Second)

	t.Logf("[%s] done", end)
	if end != 17*time.Second {
		t.Fatalf("expected EOF at 9 seconds; actual %s", end)
	}
}
