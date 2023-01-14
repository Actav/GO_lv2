package main

import (
  "fmt"
  "time"
)

var (
	a, b int
)

type ErrorWithTrace struct {
  err interface{}
  time string
}

func New(err interface{}) error {
  // получим текущее время
    dt := time.Now()

  return &ErrorWithTrace{
    err: err,                                      // ошибкa
    time: dt.Format("2006.01.02 15:04:05.000000"), // кастуем время в нужный формат
  }
}
func (e *ErrorWithTrace) Error() string {
  return fmt.Sprintf("\n  error: %s\n  time: %s", e.err, e.time)
}


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
          fmt.Println("recovered ->", New(v))
      }
    }()

  // Вывыдем результат деления двух чисел
    fmt.Printf("Результат деления введенных чисел: %d", a/b)
}
