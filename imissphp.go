package imissphp

// A pacakge for common functions that I cannot find in the standard library

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

// Recursively converts a value of type `any` into a map[string]... structure.
// Slices and arrays are converted to map[string]... as well with the index used as a string key.
// Structs are converted using json.Marshal so that json struct tags are used.
// If value is not a map, struct, slice, or array then an empty map `map[string]any{}` is returned
func ToMap(value any) map[string]any {
	// Get the type of the value
	val := reflect.ValueOf(value)
	processableKinds := []reflect.Kind{reflect.Map, reflect.Struct, reflect.Array, reflect.Slice}

	if val.Kind() == reflect.Map {
		// For maps, we only need to recurse for nested values
		result := make(map[string]any)

		for _, key := range val.MapKeys() {
			// Get the value at the current key
			mapValue := val.MapIndex(key)
			isProcessable := InArray(mapValue.Kind(), processableKinds)

			// Recursively convert if map or struct
			if isProcessable {
				nestedValue := ToMap(mapValue.Interface())
				result[key.String()] = nestedValue
			} else {
				// Otherwise, directly assign the value
				result[key.String()] = mapValue.Interface()
			}
		}

		return result
	} else if val.Kind() == reflect.Struct {
		// Handle struct: use json.Marshal to get the struct as a map using json tags
		result := make(map[string]any)

		data, err := json.Marshal(value)
		if err != nil {
			return map[string]any{}
		}

		// Unmarshal the JSON data into a map
		if err := json.Unmarshal(data, &result); err != nil {
			return map[string]any{}
		}

		// Now recurse
		for key, val := range result {
			// Recurse on nested structures
			mapVal := reflect.ValueOf(val)

			isProcessable := InArray(mapVal.Kind(), processableKinds)

			if isProcessable {
				nestedValue := ToMap(val)
				result[key] = nestedValue
			} else {
				result[key] = val
			}
		}

		return result
	} else if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		// Handle slices and arrays: iterate over elements, convert to map using string index as keys
		result := make(map[string]any)

		for i := 0; i < val.Len(); i++ {
			element := val.Index(i)

			// Convert the index to string and recursively convert each element
			key := strconv.Itoa(i)

			isProcessable := InArray(element.Kind(), processableKinds)

			if isProcessable {
				nestedValue := ToMap(element.Interface())
				result[key] = nestedValue
			} else {
				result[key] = element.Interface()
			}

		}

		return result
	} else {
		// The value passed in cannot be processed, return empty map
		return map[string]any{}
	}
}

// Flatten a map[string]any
func FlattenMap(node map[string]any) map[string]any {
	var _flattenMap func(node map[string]any, key string) map[string]any

	// Use an internal function for recursion to avoid passing in an empty string to outer function
	_flattenMap = func(node map[string]any, key string) map[string]any {
		flattened := map[string]any{}

		// Iterate over map keys/values
		for k, v := range node {
			mapV, isMap := v.(map[string]any)
			isLeaf := !isMap

			flattenedKey := fmt.Sprintf("%s.%s", key, k)
			if key == "" {
				flattenedKey = k
			}

			if isLeaf {
				flattened[flattenedKey] = v
			} else {
				subFlattened := _flattenMap(mapV, flattenedKey)
				for _k, _v := range subFlattened {
					flattened[_k] = _v
				}
			}
		}

		return flattened
	}

	return _flattenMap(node, "")
}

// Unflatten a map[string]any
func UnFlattenMap(node map[string]any) map[string]any {
	result := map[string]any{}

	for k, v := range node {
		keys := strings.Split(k, ".")

		curr := result

		for i, key := range keys {
			if i >= len(keys)-1 {
				curr[key] = v
			} else {
				_, exists := curr[key]
				if !exists {
					curr[key] = map[string]any{}
				}
				curr, _ = curr[key].(map[string]any)
			}
		}
	}

	return result
}
