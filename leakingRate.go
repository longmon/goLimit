package goLimit

import (
	"errors"
	"sync"
	"time"
)

// defaultRateLimit 默认速率
const defaultRateLimit uint64 = 100

type leaking struct {
	sync.RWMutex               //同步锁
	rate         time.Duration //每个请求花费的时间
	last         time.Time     //上次请求发生的时间
	sleep        bool          //等待返回 depreated
	err          error         //错误类型， 因为只有一种错误类型，所以在创建的时候初始了，免得每次Take都创建
	backlog
}

// NewLeaking 初始化一个漏桶
// 传入速率参数rate 表示一秒钟释出多少水滴
// 默认速率限制时间内需要等待
func NewLeaking(rate uint64) *leaking {
	if rate == 0 {
		rate = defaultRateLimit
	}
	l := new(leaking)
	l.rate = time.Second / time.Duration(rate)
	l.last = time.Now()
	l.err = errors.New("rate limited")
	return l
}

// Wait 方法以恒定速率产生time.Time并返回，单位时间速率内多个请求会阻塞
func (l *leaking) Wait() *time.Time {
	l.Lock()
	defer l.Unlock()
	now := time.Now()
	//sleep 有可能是负数，因为两次请求的时间间隔有可能大于速率
	sleep := l.rate - now.Sub(l.last)
	if sleep > 0 {
		time.Sleep(sleep)
		l.last = now.Add(sleep)
	} else {
		l.last = now
	}

	return &l.last
}

// Wait 方法以恒定速率产生time.Time并返回
// 单位时间速率内多个请求不会阻塞,马上返回 error
func (l *leaking) Take() (*time.Time, error) {
	l.Lock()
	defer l.Unlock()
	now := time.Now()
	sleep := l.rate - now.Sub(l.last)
	if sleep > 0 {
		return nil, l.err
	}
	l.last = now.Add(sleep)
	return &l.last, nil
}
