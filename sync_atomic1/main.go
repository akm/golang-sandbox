package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	countUp := NewAtomicCountup()
	ctx, cancel := context.WithCancel(context.Background())
	go watch(ctx, countUp)
	countUp.Do()
	cancel()
	fmt.Printf("AtomicCountup done\n")
}

func watch(ctx context.Context, countUp *AtomicCountup) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		waiting, working := countUp.Stat()
		fmt.Printf("waiting: %d, working: %d\n", waiting, working)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// continue
		}
	}
}

type AtomicCountup struct {
	Waiting int64
	Working int64
	wg      sync.WaitGroup
}

func NewAtomicCountup() *AtomicCountup {
	return &AtomicCountup{}
}

func (c *AtomicCountup) Do() {
	for i := 0; i < 1000; i++ {
		c.wg.Add(1)
		atomic.AddInt64(&c.Waiting, 1)
		go func() {
			defer func() {
				c.wg.Done()
				atomic.AddInt64(&c.Working, -1)
			}()
			atomic.AddInt64(&c.Waiting, -1)
			atomic.AddInt64(&c.Working, 1)
			time.Sleep(5000 * time.Millisecond)
		}()
		time.Sleep(10 * time.Millisecond)
	}
	c.wg.Wait()
}

func (c *AtomicCountup) Stat() (int64, int64) {
	return atomic.LoadInt64(&c.Waiting), atomic.LoadInt64(&c.Working)
}
