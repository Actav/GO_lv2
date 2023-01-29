/*
С помощью пула воркеров написать программу, которая запускает 1000 горутин,
каждая из которых увеличивает число на 1. Дождаться завершения всех горутин и
убедиться, что при каждом запуске программы итоговое число равно 1000.
*/

package main

import (
  "fmt"
  "time"
)

func main() {
  count := 0
  c := make(chan int, 10)
  defer func ()  {
    fmt.Println("count ==", count)

  }()

  go func() {
    defer func ()  {
      <-c
    }()

    var i int
    for {
      c <- i
      i++
    }
  }()

  go func() {
    for f := range c {
      if f > 999 {
        return
      }
      count++
    }
  }()

  // select{}
  time.Sleep(2 * time.Second)
}
