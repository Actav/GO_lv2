/*
Написать программу, которая при получении в канал сигнала SIGTERM
останавливается не позднее, чем за одну секунду (установить таймаут).
*/

package main

import (
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "time"
)

func main() {
  // Создаем канал для получения сигналов
  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, syscall.SIGTERM)

  // Ожидаем получение сигнала SIGTERM
  <-sigs

  // Выводим сообщение о начале завершения
  fmt.Println("Exiting... Please wait.")

  // Устанавливаем таймер для завершения программы
  timeout := time.After(time.Second)

  select {
    case <-timeout:{
      fmt.Println("Exited.")
      os.Exit(0)
    }
  }
}
