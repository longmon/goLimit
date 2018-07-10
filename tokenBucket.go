package goLimit

import (
	"errors"
	"sync"
	"time"
)

var (
	defaultbSize    int    = 100
	defaultbGenRate uint64 = 100
)

// Bucket 令牌桶
type Bucket struct {
	sync.RWMutex
	chant chan time.Time
	leaky *leaking
	rate  uint64
	err   error
	wait  time.Duration
}

// NewBucket 初始化一个令牌桶实现
func NewBucket(size int, rate uint64, waitfor time.Duration) *Bucket {
	if size <= 0 {
		size = defaultbSize
	}
	if rate == 0 {
		rate = defaultbGenRate
	}
	b := new(Bucket)
	b.chant = make(chan time.Time, size)
	b.leaky = NewLeaking(rate)
	b.err = errors.New("Empty bucket ( or time out )!")
	b.rate = rate
	b.wait = waitfor
	go b.produce()
	return b
}

//以恒定速率往桶里面放令牌
func (b *Bucket) produce() {
	b.Lock()
	defer b.Unlock()
	for {
		t := b.leaky.Wait()
		select {
		case b.chant <- *t:
		}
	}
}

// Get 从令牌桶里面取出一个令牌，不阻塞，如果令牌桶为空，返回nil,b.err
func (b *Bucket) Take() (time.Time, error) {
	select {
	case t := <-b.chant:
		return t, nil
	default:
		return time.Time{}, b.err
	}
}

// Wait 阻塞等待令牌,直到超时
func (b *Bucket) Wait() (time.Time, error) {
	if b.wait == 0 {
		select {
		case t := <-b.chant:
			return t, nil
		}
	} else {
		select {
		case <-time.After(b.wait):
			return time.Time{}, b.err
		case t := <-b.chant:
			return t, nil
		}
	}

	return time.Time{}, b.err
}

// Info 返回令牌桶令牌
func (b *Bucket) Info() map[string]int {
	info := make(map[string]int, 3)
	info["len"] = len(b.chant)
	info["cap"] = cap(b.chant)
	info["rate"] = int(b.rate)
	return info
}
