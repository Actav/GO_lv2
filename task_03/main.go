package main

import (
  "errors"
  "fmt"
  "os"
)
func main() {
  // отложенная функця обработки panic()
    defer func() {
      if v := recover(); v != nil {
        fmt.Println("recovered", v)
      }
    }()

  // проверим есть ли папка назначения, если нет создадим
    if _, err := os.Stat("./emptyDir"); errors.Is(err, os.ErrNotExist) {
      err := os.Mkdir("./emptyDir", 0777)
      if err != nil {
        panic(err)
      }
    }

  // сгенерируем в цикле N пустых файлов
    for i := 0; i < 10; i++ {
      // создаим новый файл согласно индексу интерации
        f, err := os.Create(fmt.Sprintf("./emptyDir/%d", i))
        if err != nil {
          panic(err)
        }

      // отложенно закроем соединение с файлом
        defer f.Close()
    }
}
