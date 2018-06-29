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

type Token struct {
	t time.Time
}

type Bucket struct {
	chant chan Token
	leaky *leaking
	rate  uint64
	err   error
	one   *sync.Once
}

func InitBucket(size int, rate uint64) *Bucket {
	if size <= 0 {
		size = defaultBucketSize
	}
	if rate == 0 {
		rate = defaultBucketGenRate
	}
	bucket := new(Bucket)
	bucket.chant = make(chan Token, size)
	bucket.leaky = NewLeaking(rate).WithSleep()
	bucket.err = errors.New("Empty Bucket!")
	bucket.rate = rate
	bucket.one = &sync.Once{}
	return bucket
}

func (bucket *Bucket) Produce() {
	bucket.one.Do(func() {
		bucket.doProduceOnce()
	})
}

func (bucket *Bucket) doProduceOnce() {
	for {
		t, _ := bucket.leaky.Take()
		select {
		case bucket.chant <- Token{t: *t}:
		}
	}
}

func (bucket *Bucket) Take() (*Token, error) {
	select {
	case token := <-bucket.chant:
		return &token, nil
	default:
		return nil, bucket.err
	}
}

func (bucket *Bucket) Wait() (*Token, error) {
	select {
	case token := <-bucket.chant:
		return &token, nil
	}
	return nil, bucket.err
}

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
