package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println(i)
			runtime.Gosched()
			fmt.Println(i + 10)
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}
