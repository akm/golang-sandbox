package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type WaitGroupWriter struct {
	io.WriteCloser
	wg *sync.WaitGroup
}

func (w *WaitGroupWriter) Close() error {
	if err := w.WriteCloser.Close(); err != nil {
		return err
	}
	w.wg.Wait()
	return nil
}

func main() {
	longBytes := multipleBytes([]byte("0123456789abcdef"), 65) // バッファサイズよりちょっとだけ大きいデータ
	// fmt.Printf("%s\n", longBytes)

	r, pw := io.Pipe()
	wg := &sync.WaitGroup{}
	w := &WaitGroupWriter{WriteCloser: pw, wg: wg}

	wg.Add(1)
	go func() {
		fmt.Printf("before io.Copy()\n")
		defer fmt.Printf("after  io.Copy()\n")
		buf := make([]byte, 1024)
		if _, err := io.CopyBuffer(os.Stdout, r, buf); err != nil {
			fmt.Printf("error  io.Copy(): %+v\n", err)
			os.Exit(1)
		}
		wg.Done()
	}()

	for i := 0; i < 3; i++ {
		fmt.Fprintf(w, "%s\n", longBytes)
		time.Sleep(time.Second)
	}
	fmt.Printf("before w.Close()\n")
	if err := w.Close(); err != nil {
		fmt.Printf("error  w.Close(): %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("after  w.Close()\n")

	// Result:
	// $ go run .
	// before io.Copy()
	// 0123456789abcdef0123456789abcdef...
	// 0123456789abcdef0123456789abcdef...
	// 0123456789abcdef0123456789abcdef...
	// before w.Close()
	// after  w.Close()
	// after  io.Copy()
}

func multipleBytes(base []byte, times int) []byte {
	r := make([]byte, len(base)*times)
	for i := 0; i < times; i++ {
		copy(r[i*len(base):], base)
	}
	return r
}
