package golimit

import (
	"time"
)

const defaultRateLimit uint64 = 100

type Leaking struct {
	rate  time.Duration //每个请求花费的时间
	last  time.Time     //上次请求发生的时间
	sleep bool
}

func NewLeaking(rate uint64) *Leaking {
	if rate == 0 {
		rate = defaultRateLimit
	}
	l := new(Leaking)
	l.rate = time.Second / time.Duration(rate)
	return nil
}

func (l *Leaking) Take() (time.Duration, error) {
	now := time.Now()
	//sleep 有可能是负数，因为两次请求的时间间隔有可能大于速率
	sleep := l.rate - now.Sub(l.last)
	if sleep > 0 {
		time.Sleep(sleep)
	}

}
