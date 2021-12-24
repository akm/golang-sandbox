package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	countUp := NewAtomicCountup()
	countUp.Do()
	fmt.Printf("AtomicCountup: %d\n", countUp.Count)
}

type AtomicCountup struct {
	Count uint32
	wg    sync.WaitGroup
}

func NewAtomicCountup() *AtomicCountup {
	return &AtomicCountup{}
}

func (c *AtomicCountup) Do() {
	for i := 0; i < 10000; i++ {
		c.wg.Add(1)
		go func() {
			atomic.AddUint32(&c.Count, 1)
			defer c.wg.Done()
		}()
	}
	c.wg.Wait()
}
