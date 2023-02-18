package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Struct represents a generated struct
type Struct struct {
	Name   string
	Fields []Field
}

// Field represents a generated field
type Field struct {
	Name string
	Type string
	Tag  string
}

// Generate generates the code for a given struct
func (s Struct) Generate() string {
	var sb strings.Builder

	// Write struct declaration
	fmt.Fprintf(&sb, "type %s struct {\n", s.Name)

	// Write fields
	for _, f := range s.Fields {
		fmt.Fprintf(&sb, "\t%s %s %s\n", f.Name, f.Type, f.Tag)
	}

	// Close struct declaration
	sb.WriteString("}")

	return sb.String()
}

// WriteToFile writes the generated code to a file
func WriteToFile(filename, contents string) error {
	return ioutil.WriteFile(filename, []byte(contents), 0644)
}

func main() {
	// Create a new struct
	myStruct := Struct{
		Name: "Person",
		Fields: []Field{
			{Name: "Name", Type: "string", Tag: "`json:\"name\"`"},
			{Name: "Age", Type: "int", Tag: "`json:\"age\"`"},
		},
	}

	// Generate code for the struct
	code := myStruct.Generate()

	// Write the code to a file
	err := WriteToFile("person.go", code)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Println("Code generated and written to file!")
}
