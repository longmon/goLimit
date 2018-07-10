package goLimit

import (
	"errors"
	"fmt"
	"sync"
	"time"
)
// Pool是一个预定义的池，在创建的时候设置池的大小与池中没数据时客户端等待的延时
// 客户端处理事务完成后，可以把数据重新放回池或丢弃
type Pool struct {
	cap  uint
	ch   chan interface{}
	wait time.Duration
	err  []error
	o    *sync.Once
}

//NewPool 得到一个池
func NewPool(size uint, waitfor time.Duration) *Pool {
	if size == 0 {
		size = 1
	}

	p := new(Pool)
	p.cap = size
	p.ch = make(chan interface{}, size)
	p.wait = waitfor
	p.err = []error{
		errors.New("empty pool"),
		errors.New(fmt.Sprintf("wait time out on empty pool")),
	}
	return p
}

//Init初始化，在获得一个池后需要进行初始化，否则池中无数据啊
//本方法多次执行无效
func (p *Pool) Init(v interface{}) {
	p.o.Do(func() {
		for {
			select {
			case p.ch <- v:
			default:
				break
			}
		}
	})
}

// Take 非阻塞取出数据
func (p *Pool) Take() (interface{}, error) {
	select {
	case v := <-p.ch:
		return v, nil
	default:
		return nil, p.err[0]
	}
}

// Wait 等待数据，直到超时
func (p *Pool) Wait() (interface{}, error) {
	select {
	case <-time.After(p.wait):
		return nil, p.err[1]
	case v := <-p.ch:
		return v, nil
	}
}

// 把数据放回池中
func (p *Pool) Put(v interface{}) {
	select {
	case p.ch <- v:
	default:
		return
	}
}
