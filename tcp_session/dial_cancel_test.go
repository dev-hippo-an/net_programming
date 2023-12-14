package tcp_session

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

// 데드라인 지정 없이 필요 시 cancel 함수를 사용하여 연결 시도 취소
func TestDialContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	sync := make(chan struct{})

	go func() {
		defer func() {
			sync <- struct{}{}
		}()

		var d net.Dialer

		d.Control = func(_, _ string, _ syscall.RawConn) error {
			time.Sleep(time.Second)
			return nil
		}

		conn, err := d.DialContext(ctx, "tcp", "10.0.0.1:80")
		if err != nil {
			t.Logf("Error occured caused by canceled context; %s", err)
			return
		}

		t.Log("Unreachable code here..!")

		conn.Close()
		t.Error("connection did not time out")
	}()

	cancel()
	<-sync
	defer close(sync)

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %q", ctx.Err())
	} else {
		t.Log("context cancel")
	}
}
