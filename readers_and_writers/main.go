package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	// bufio.NewWriter のデフォルトのバッファサイズは 4096 bytes
	// https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/bufio/bufio.go;drc=ceda93ed673294f0ce5eb3a723d563091bff0a39;l=19
	longBytes := multipleBytes([]byte("0123456789abcdef"), 65) // バッファサイズよりちょっとだけ大きいデータ
	// fmt.Printf("%s\n", longBytes)

	r, w := io.Pipe()

	go func() {
		for i := 0; i < 3; i++ {
			fmt.Fprintf(w, "%s\n", longBytes)
			time.Sleep(time.Second)
		}
		fmt.Printf("before w.Close()\n")
		w.Close()
		fmt.Printf("after  w.Close()\n")
	}()

	fmt.Printf("before io.Copy()\n")
	buf := make([]byte, 1024)
	if _, err := io.CopyBuffer(os.Stdout, r, buf); err != nil {
		fmt.Printf("error  io.Copy(): %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("after  io.Copy()\n")

}

func multipleBytes(base []byte, times int) []byte {
	r := make([]byte, len(base)*times)
	for i := 0; i < times; i++ {
		copy(r[i*len(base):], base)
	}
	return r
}
