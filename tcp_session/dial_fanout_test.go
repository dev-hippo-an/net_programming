package tcp_session

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestDialContextCancelFanOut(t *testing.T) {

	// 10초의 데드라인과 함께 컨텍스트와 캔슬 함수 반환
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(10*time.Second),
	)

	// 리스너 설정
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// 리스너 자원 해제
	defer func() {
		t.Log("메인 고루틴에서 마지막입니다 : 자원 정리 시간....")
		listener.Close()
	}()

	// 하나의 연결만 수락하고 연결 종료
	go func() {
		t.Log("Go routine started!")
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
		t.Log("Go routine ended!")
	}()

	dial := func(ctx context.Context, address string, response chan int, id int, wg *sync.WaitGroup) {
		defer func() {
			t.Logf("Dial 함수 고루틴 안에서 마지막입니다 : 자원 정리 시간....")
			wg.Done()
		}()

		var d net.Dialer
		c, err := d.DialContext(ctx, "tcp", address)
		if err != nil {
			t.Logf("Error 가 발생했습니다. 아마 Connection closed error : %v", err)
			return
		}

		c.Close()

		select {
		case <-ctx.Done():
		case response <- id:
		}
	}

	res := make(chan int)
	var wg sync.WaitGroup

	// wait group 을 이용하여 여러 고루틴 자원 제어
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go dial(ctx, listener.Addr().String(), res, i+1, &wg)
	}

	response := <-res // res 채널에서 데이터 수신까지 블로킹
	cancel()
	wg.Wait()
	close(res)

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %s", ctx.Err())
	}

	t.Logf("dialer %d retrieved the resource", response)

}
