package main

import (
	"fmt"
	"sync"
)

var counter int
var wg sync.WaitGroup

func incrementCounter(id int) {
	defer wg.Done()
	counter++
	fmt.Printf("Counter value after goroutine %d: %d\n", id, counter)
}

func main() {
	wg.Add(2)
	go incrementCounter(1)
	go incrementCounter(2)
	wg.Wait()
	fmt.Println("Final Counter:", counter)
}
