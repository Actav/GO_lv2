package main

import (
	"errors"
	"reflect"
	"testing"
)

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

type TestStruct struct {
	Foo int
	Bar string
}

func TestSetFieldsByMap(t *testing.T) {
	testCases := []struct {
		name     string
		in       *TestStruct
		values   map[string]interface{}
		expected *TestStruct
		wantErr  bool
	}{
		{
			name: "set one field",
			in:   &TestStruct{},
			values: map[string]interface{}{
				"Foo": 42,
			},
			expected: &TestStruct{
				Foo: 42,
			},
			wantErr: false,
		},
		{
			name: "set multiple fields",
			in:   &TestStruct{},
			values: map[string]interface{}{
				"Foo": 42,
				"Bar": "hello",
			},
			expected: &TestStruct{
				Foo: 42,
				Bar: "hello",
			},
			wantErr: false,
		},
		{
			name: "unknown field",
			in:   &TestStruct{},
			values: map[string]interface{}{
				"Baz": true,
			},
			expected: &TestStruct{},
			wantErr:  true,
		},
		{
			name: "field cannot be set",
			in: &TestStruct{
				Foo: 42,
			},
			values: map[string]interface{}{
				"Foo": 43,
			},
			expected: &TestStruct{
				Foo: 42,
			},
			wantErr: true,
		},
		{
			name: "type mismatch",
			in: &TestStruct{
				Foo: 42,
			},
			values: map[string]interface{}{
				"Foo": "hello",
			},
			expected: &TestStruct{
				Foo: 42,
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetFieldsByMap(tc.in, tc.values)
			if (err != nil) != tc.wantErr {
				t.Errorf("got error %v, want error %v", err, tc.wantErr)
			}
			if !reflect.DeepEqual(tc.in, tc.expected) {
				t.Errorf("got %v, want %v", tc.in, tc.expected)
			}
		})
	}
}
