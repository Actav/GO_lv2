/*
С помощью пула воркеров написать программу, которая запускает 1000 горутин,
каждая из которых увеличивает число на 1. Дождаться завершения всех горутин и
убедиться, что при каждом запуске программы итоговое число равно 1000.
*/

package main

import (
  "fmt"
  "sync"
)

var (
  wg sync.WaitGroup
  m sync.Mutex
  sem = make(chan struct{}, 50)
)

func main() {
  counter := 0

  for i := 0; i < 1000; i++ {
    wg.Add(1)
    sem <- struct{}{}

    go func() {
      m.Lock()
        counter++
      m.Unlock()

      <-sem
      wg.Done()
    }()
  }

  wg.Wait()
  fmt.Println("Counter:", counter)
}
