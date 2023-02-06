// Реализуйте функцию для разблокировки мьютекса с помощью defer

package main

import (
  "fmt"
  "sync"
)

var (
  wg sync.WaitGroup
  m sync.Mutex
)

func main() {
  counter := 0

  for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
      m.Lock()
      defer m.Unlock()

      counter++
      wg.Done()
    }()
  }

  wg.Wait()
  fmt.Println("Counter:", counter)
}
