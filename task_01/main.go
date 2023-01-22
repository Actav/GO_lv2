
package main

import (
  "div"

  "fmt"
)

func main() {
  var a, b int

  // получим строку с входными данными с утсройства ввода
    fmt.Print("Введите целое число: ")
    fmt.Scanln(&a)
    fmt.Print("Введите целое число: ")
    fmt.Scanln(&b)

  res := div.Div(a, b)

  // Вывыдем результат деления двух чисел
    fmt.Printf("Результат деления введенных чисел: %d", res)
}
