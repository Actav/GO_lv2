/*
Напишите программу, которая запускает потоков и дожидается завершения их всех
*/

package main

import (
  "sync"
)

var (
  count = 999
  wg sync.WaitGroup
  ch = make(chan struct{}, count)
)

func main() {
  for i := 0; i < count; i++ {
    wg.Add(1)
    go func() {
      ch <- struct{}{}
      wg.Done()
    }()
  }

  wg.Wait()
  close(ch)

  i := 0
  for range ch {
    i++
  }
  println(i)
}
