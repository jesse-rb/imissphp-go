package imissphp

// A pacakge for common functions that I cannot find in the standard library

import (
	"encoding/json"
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

// Gets the name of the supplied value's type.
// The i parameter can be either a value or a pointer to a value.
// This function returns the string name of the supplied value's type.
func TypeName(i interface{}) string {
	data := reflect.TypeOf(i)

	data = normalizeReflectType(data)

	return data.Name()
}

// MethodExists checks if a method with the given name exists on a type.
// The i parameter can be either a value or a pointer to a value.
// The methodName can refer to a method with either a value receiver or a pointer receiver.
// The function returns true if the method exists on either the value or its pointer type.
func MethodExists(i interface{}, methodName string) bool {
	data := reflect.TypeOf(i)

	data = normalizeReflectType(data)

	// Checks for MyStruct{}.MyMethod OR *MyStruct{}.MyMethod
	_, hasMethod := data.MethodByName(methodName)
	_, ptrHasMethod := reflect.PointerTo(data).MethodByName(methodName)

	return hasMethod || ptrHasMethod
}

// Converts any into a map[string]any.
// Useful to convert stucts into map[string]any, e.g. when using the gin.GinH type from the https://github.com/gin-gonic/gin framework.
// If the function fails to make the conversion, an empty map[string]any{} value is returend.
// Returns data converted to map[string]any by calling json.Marshal on data, so that json struct tags are respected,
// then unmarshalling into a map[string]any
func ToMap(data any) map[string]any {
	// Marshal the struct into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return map[string]any{}
	}

	// Unmarshal the JSON into a map
	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return map[string]any{}
	}

	return result
}

// func FlattenMap(map[string]any) map[string]any {
// }
//
// func UnFlattenMap(map[string]any) map[string]any {
// }
