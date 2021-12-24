package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	countUp := NewAtomicCountup()
	countUp.Do()
	fmt.Printf("AtomicCountup done\n")
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
