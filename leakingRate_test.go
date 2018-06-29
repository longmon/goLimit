package goLimit

import (
	"testing"
	"time"
)

func TestTake(t *testing.T) {
	leaky := NewLeaking(100)
	last := time.Now()
	for i := 0; i < 10; i++ {
		tm := leaky.Wait()
		sub := tm.Sub(last)
		last = *tm
		if sub < (time.Second / 100) {
			t.Errorf("too quickly: %v\n", sub)
		}
		if sub > time.Second/80 {
			t.Errorf("too slowly: %v\n", sub)
		}
	}
	t.Log("passed")
}

func TestTakeWithoutSleep(t *testing.T) {
	leaky := NewLeaking(100)
	last := time.Now()
	for i := 0; i < 1000; i++ {
		tm, err := leaky.Take()
		if err != nil {
			continue
		}
		sub := tm.Sub(last)
		last = *tm
		if sub < time.Second/100 {
			t.Fatalf("to quick:%v\n", sub)
		}
		if sub > time.Second/80 {
			t.Fatalf("to slowly:%v\n", sub)
		}
	}
	t.Log("passed")
}
