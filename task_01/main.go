package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"sync"
)

var data int
var lock sync.Mutex

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

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
      lock.Lock()
      defer func ()  {
        lock.Unlock()
        wg.Done()
      }()

			data++
			fmt.Println(data)
		}()
	}
  
	wg.Wait()
}
