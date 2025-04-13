package imissphpgo

// A pacakge for common functions that I cannot find in the standard library

import (
	"reflect"
	"unicode"
)

// Capitalizes the first letter in a string
func UcFirst(s string) string {
	if len(s) <= 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

// Check if item is in array
func InArray[T comparable](val T, list []T) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}

	return false
}

// Normalize a reflect.Type by returning the element type if it is a pointer
func normalizeReflectType(data reflect.Type) reflect.Type {
	if data.Kind() == reflect.Pointer {
		data = data.Elem()
	}

	return data
}

func TypeName(i interface{}) string {
	data := reflect.TypeOf(i)

	data = normalizeReflectType(data)

	return data.Name()
}

func MethodExists(i interface{}, methodName string) bool {
	data := reflect.TypeOf(i)

	data = normalizeReflectType(data)

	// Checks for MyStruct{}.MyMethod OR *MyStruct{}.MyMethod
	_, hasMethod := data.MethodByName(methodName)
	_, ptrHasMethod := reflect.PointerTo(data).MethodByName(methodName)

	return hasMethod || ptrHasMethod
}
