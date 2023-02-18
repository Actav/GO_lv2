package main

import (
	"errors"
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	person := Person{
		Name: "Alice",
		Age:  25,
	}

	values := map[string]interface{}{
		"Name": "Bob",
		"Age":  30,
	}

	if err := SetFieldsByMap(&person, values); err != nil {
		fmt.Printf("Error setting field %+v: %s\n", values, err)
	}

	fmt.Printf("Person: %+v\n", person)
}

func SetFieldsByMap(in interface{}, values map[string]interface{}) error {
	v := reflect.ValueOf(in).Elem()
	for key, val := range values {
		fieldValue := v.FieldByName(key)
		if !fieldValue.IsValid() {
			return errors.New("field " + key + " not found")
		}
		if !fieldValue.CanSet() {
			return errors.New("field " + key + " cannot be set")
		}
		fieldType := fieldValue.Type()
		valType := reflect.ValueOf(val).Type()
		if fieldType != valType {
			return errors.New("type mismatch for field " + key)
		}
		fieldValue.Set(reflect.ValueOf(val))
	}
	return nil
}
