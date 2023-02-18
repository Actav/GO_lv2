package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// CountAsyncFuncCalls подсчитывает количество вызовов асинхронных функций в заданном файле Go
// и возвращает его вместе с ошибкой
func CountAsyncFuncCalls(fileName, funcName string) (int, error) {
	// Открываем файл и парсим его
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return 0, err
	}

	// С помощью поиска в глубину обходим дерево AST и ищем вызовы целевой функции
	var count int
	ast.Inspect(node, func(n ast.Node) bool {
		// Ищем вызов целевой функции
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if funcIdent, ok := callExpr.Fun.(*ast.Ident); ok {
				if funcIdent.Name == funcName {
					// Ищем асинхронные вызовы функций
					if isAsyncCall(callExpr, node) {
						count++
					}
				}
			}
		}
		return true
	})

	return count, nil
}

// isAsyncCall возвращает true, если данный вызов является асинхронным, т.е. содержит go ключевое слово
func isAsyncCall(callExpr *ast.CallExpr, node ast.Node) bool {
	// Ищем выражение "go"
	var isAsync bool
	ast.Inspect(callExpr, func(n ast.Node) bool {
		if unaryExpr, ok := n.(*ast.UnaryExpr); ok {
			if unaryExpr.Op == token.GO {
				isAsync = true
			}
		}
		return !isAsync // прекращаем обход, если нашли "go"
	})

	// Ищем обработчики каналов внутри вызова функции
	if isAsync {
		var isChannelHandler bool
		ast.Inspect(node, func(n ast.Node) bool {
			if callExpr, ok := n.(*ast.CallExpr); ok {
				if funcIdent, ok := callExpr.Fun.(*ast.Ident); ok {
					if strings.HasSuffix(funcIdent.Name, "Chan") {
						isChannelHandler = true
					}
				}
			}
			return !isChannelHandler // прекращаем обход, если нашли обработчик канала
		})
		return !isChannelHandler
	}

	return false
}

func main() {
	count, err := CountAsyncFuncCalls("example.go", "SomeFunction")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Number of async calls: %d\n", count)
	}
}
