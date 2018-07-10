package goLimit

import (
	"errors"
	"sync"
	"time"
)

var (
	defaultBucketSize    int    = 100
	defaultBucketGenRate uint64 = 100
)

// Token 令牌
type Token struct {
	t time.Time
}

// Bucket 令牌桶
type Bucket struct {
	sync.RWMutex
	chant chan Token
	leaky *leaking
	rate  uint64
	err   error
}

// InitBucket 初始化一个令牌桶实现
func InitBucket(size int, rate uint64) *Bucket {
	if size <= 0 {
		size = defaultBucketSize
	}
	if rate == 0 {
		rate = defaultBucketGenRate
	}
	bucket := new(Bucket)
	bucket.chant = make(chan Token, size)
	bucket.leaky = NewLeaking(rate)
	bucket.err = errors.New("Empty Bucket!")
	bucket.rate = rate
	go bucket.produce()
	return bucket
}

//以恒定速率往桶里面放令牌
func (bucket *Bucket) produce() {
	bucket.Lock()
	defer bucket.Unlock()
	for {
		t, _ := bucket.leaky.Wait()
		select {
		case bucket.chant <- Token{t: *t}:
		}
	}
}

// Get 从令牌桶里面取出一个令牌，不阻塞，如果令牌桶为空，返回nil,bucket.err
func (bucket *Bucket) Take() (*Token, error) {
	select {
	case token := <-bucket.chant:
		return &token, nil
	default:
		return nil, bucket.err
	}
}

// Wait 阻塞等待令牌
func (bucket *Bucket) Wait() (*Token, error) {
	select {
	case token := <-bucket.chant:
		return &token, nil
	}
	return nil, bucket.err
}

// Info 返回令牌桶令牌
func (bucket *Bucket) Info() map[string]int {
	info := make(map[string]int, 3)
	info["len"] = len(bucket.chant)
	info["cap"] = cap(bucket.chant)
	info["rate"] = int(bucket.rate)
	return info
}

func (t *Token) GetTokenTime() time.Time {
	return t.t
}
