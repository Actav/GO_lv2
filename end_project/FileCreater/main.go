package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Установить seed для случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Получить путь к директории из флага командной строки
	dirPathPtr := flag.String("dir", "./test_dir", "Путь к директории")
	flag.Parse()

	if *dirPathPtr == "" {
		fmt.Println("Пожалуйста, укажите путь к директории через флаг '-dir'")
		os.Exit(1)
	}

	// Пройти по всем файлам в директории
	err := filepath.Walk(*dirPathPtr, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропустить если это директория
		if info.IsDir() {
			return nil
		}

		// Создать копию файла
		err = copyFile(path, *dirPathPtr)
		if err != nil {
			return err
		}

		fmt.Printf("Создана копия файла %s\n", path)

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

// Функция для копирования файла
func copyFile(src string, destDir string) error {
	// Открыть исходный файл для чтения
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Получить имя файла без пути
	filename := filepath.Base(src)

	// Создать новую случайную директорию внутри директории назначения
	dest, err := createRandomDirectory(destDir)
	if err != nil {
		return err
	}

	// Создать новый файл в новой директории
	dest = filepath.Join(dest, filename)
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Скопировать содержимое исходного файла в новый файл
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// Функция создания новой случайной директории внутри директории назначения
func createRandomDirectory(destDir string) (string, error) {
	// Создать новую случайную директорию внутри директории назначения
	randDir := fmt.Sprintf("%d", rand.Intn(1000))
	dest := filepath.Join(destDir, randDir)
	err := os.Mkdir(dest, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dest, nil
}
