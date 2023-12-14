package tcp_session

import (
	"context"
	"fmt"
	"io"
	"time"
)

const defaultPingInterval = 30 * time.Second

func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	var interval time.Duration
	select {
	case <-ctx.Done():
		return
	case interval = <-reset: // 1 reset 채널에서 초기 간격을 받아 옴 - reset 의 초기 cap 이 1 이기 때문에 블로킹 없이 값을 받아옮
		fmt.Println("this is initial interval from channel : ", interval)
	default: // 초기 값이 제공되지 않는 경우 default 로 흐름 진행
		fmt.Println("this is default select statement.. let's get it")
	}

	if interval <= 0 {
		interval = defaultPingInterval
	}

	timer := time.NewTimer(interval) // 2

	defer func() {
		if !timer.Stop() {
			time := <-timer.C
			fmt.Println("this is time : ", time)
		}
	}()

	for { // 3가지 케이스 중 하나가 일어날 때까지 블로킹
		select {
		case <-ctx.Done(): // 3 컨텍스트가 취소된 경우
			return
		case newInterval := <-reset: // 4 리셋을 위한 시그널을 받은 경우 인터벌 리셋
			if !timer.Stop() {
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C: // 5 타이머 만료
			if _, err := w.Write([]byte("ping")); err != nil {
				return
			}
		}
		// 이후 다시 루프를 시작하기 전에 interval 로 timer reset 함
		_ = timer.Reset(interval) // 6
	}
}
