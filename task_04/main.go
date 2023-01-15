package main

import (
  "fmt"
  "time"
)

func main() {
  // запустим отдельную горутину
    go func() {
      // отложенная функця обработки panic()
        defer func() {
          if v := recover(); v != nil {
            fmt.Println("recovered", v)
          }
        }()

      // явный вызов паники
        panic("A-A-A!!!")
    }()

  // пауза выполнения
    time.Sleep(time.Second * 10)
}
