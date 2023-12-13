package tcp_session

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

// 명시적으로 타임 아웃을 정의하기 with context and cancel
func TestDialContext(t *testing.T) {
	// deadline 설정
	dl := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), dl)
	defer cancel()

	var d net.Dialer
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		time.Sleep(5*time.Second + time.Millisecond)
		return nil
	}

	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
	if err == nil {
		conn.Close()
		t.Fatal("connection did not time out")
	}

	nErr, ok := err.(net.Error)

	if !ok {
		t.Error(err)
	} else {
		if !nErr.Timeout() {
			t.Errorf("error is not a timeout: %v", err)
		}
	}

	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}
	// 결국 context.DeadlineExceeded 가 발생하고 defer 키워드로 cancel() 함수가 호출되어 graceful terminate 
}
