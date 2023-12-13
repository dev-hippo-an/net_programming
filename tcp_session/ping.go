package tcp_session

import (
	"context"
	"io"
	"time"
)

const defaultPingInterval = 30 * time.Second

func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	var interval time.Duration
	select {
	case <- ctx.Done():
		return
	case interval = <-reset:  // 1 reset 채널에서 초기 간격을 받아 옴
	default:
	}

	if interval <= 0 {
		interval = defaultPingInterval
	}

	timer := time.NewTimer(interval)  // 2

	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {
		case <-ctx.Done():  // 3
			return
		case newInterval := <-reset:  // 4
			if !timer.Stop() {
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C:  // 5
			if _, err := w.Write([]byte("ping")); err != nil {
				return
			}
		}

		_ = timer.Reset(interval)  // 6
	}
}
