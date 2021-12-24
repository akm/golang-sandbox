package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	cnt := useAtomicAddUint32()
	fmt.Printf("useAtomicAddUint32: %d\n", cnt)
}

func useAtomicAddUint32() uint32 {
	var cnt uint32
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			atomic.AddUint32(&cnt, 1)
			defer wg.Done()
		}()
	}

	wg.Wait()

	return cnt
}
