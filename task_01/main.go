package main

import (
  "fmt"
)

var (
	a, b int
)

func main() {
  // получим строку с входными данными с утсройства ввода
    fmt.Print("Введите целое число: ")
    fmt.Scanln(&a)
    fmt.Print("Введите целое число: ")
    fmt.Scanln(&b)

  // Отлженный вызов оброботчика panic()
    defer func() {
      if v := recover(); v != nil {
        // выведем в консоль полученную ошибку вызвавшую панику
          fmt.Println("recovered ->", v)
      }
    }()

  // Вывыдем результат деления двух чисел
    fmt.Printf("Результат деления введенных чисел: %d", a/b)
}
