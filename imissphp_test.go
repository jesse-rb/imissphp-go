package imissphp

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestUcFirst(t *testing.T) {
	test := "this is a test."
	expected := "This is a test."
	actual := UcFirst(test)

	if actual != expected {
		t.Fatalf("expected UcFirst did not match actual UcFirst output.")
	}
}

func TestInArray(t *testing.T) {
	testStringArray := []string{"red", "green", "blue"}

	if InArray("red", testStringArray) == false {
		t.Fatalf("Expected red to be in the array")
	}

	if InArray("yellow", testStringArray) == true {
		t.Fatalf("Did not expect yellow to be in the array")
	}

	testIntArray := []int{3, -2}

	if InArray(-2, testIntArray) == false {
		t.Fatalf("Expected -2 to be in the array")
	}

	if InArray(0, testIntArray) == true {
		t.Fatalf("Did not expect 0 to be in the array")
	}
}

func TestTypeName(t *testing.T) {
	expectedA := "Logger"
	testA := &log.Logger{}
	actualA := TypeName(testA)

	if actualA != expectedA {
		t.Fatalf("Expected struct %s to have name %s", actualA, expectedA)
	}

	expectedB := "Time"
	testB := time.Time{}
	actualB := TypeName(testB)

	if actualB != expectedB {
		t.Fatalf("Expected struct %s to have name %s", actualB, expectedB)
	}
}

func TestMethodExists(t *testing.T) {
	testStructA := log.Logger{}
	testMethodA1 := "Fatal"
	testMethodA2 := "DebugzzzzNotReal"

	if MethodExists(&testStructA, testMethodA1) == false {
		t.Fatalf("Expected method %s to exist on struct %s", testMethodA1, TypeName(&testStructA))
	}

	if MethodExists(&testStructA, testMethodA2) == true {
		t.Fatalf("Did not expect method %s to exist on struct %s", testMethodA2, TypeName(&testStructA))
	}

	testStructB := &time.Time{}
	testMethodB1 := "Date"
	testMethodB2 := "NZTimezzzzDefinitelyReal"

	if MethodExists(testStructB, testMethodB1) == false {
		t.Fatalf("Expected method %s to exist on struct %s", testMethodB1, TypeName(testStructB))
	}

	if MethodExists(testStructB, testMethodB2) == true {
		t.Fatalf("Did not expect method %s to exist on struct %s", testMethodB2, TypeName(testStructB))
	}
}

func TestToMap(t *testing.T) {
	type structA struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Words     []string
		isPrivate bool
	}

	inA := structA{
		ID:        3,
		Name:      "John",
		Words:     []string{"Hello", "there"},
		isPrivate: true,
	}

	expectedA := map[string]any{
		"id":    3,
		"name":  "John",
		"Words": map[string]string{"0": "Hello", "1": "there"},
	}

	expectedJsonA, err := json.Marshal(expectedA)
	if err != nil {
		t.Fatalf("Error converting expected map to json: %#v", err)
	}

	testJsonA, err := json.Marshal(ToMap(inA))
	if err != nil {
		t.Fatalf("Error converting test map to json: %#v", err)
	}

	if string(expectedJsonA) != string(testJsonA) {
		t.Fatalf("Expected JSON: %#v, got JSON: %#v", string(expectedJsonA), string(testJsonA))
	}
}

func TestFlattenMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		prefix   string
		expected map[string]any
	}{
		{
			name: "simple nested map",
			input: map[string]any{
				"user": map[string]any{
					"id":   3,
					"name": "john",
				},
			},
			prefix: "",
			expected: map[string]any{
				"user.id":   3,
				"user.name": "john",
			},
		},
		{
			name: "nested map with prefix",
			input: map[string]any{
				"user": map[string]any{
					"details": map[string]any{
						"age":  30,
						"city": "Paris",
					},
				},
			},
			prefix: "",
			expected: map[string]any{
				"user.details.age":  30,
				"user.details.city": "Paris",
			},
		},
		{
			name: "flat map with prefix",
			input: map[string]any{
				"id":   1,
				"name": "Alice",
			},
			prefix: "profile",
			expected: map[string]any{
				"profile.id":   1,
				"profile.name": "Alice",
			},
		},
		{
			name: "deeply nested map",
			input: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": map[string]any{
							"d": 42,
						},
					},
				},
			},
			prefix: "",
			expected: map[string]any{
				"a.b.c.d": 42,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := FlattenMap(tc.input, tc.prefix)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("FlattenMap() = %v, expected %v", got, tc.expected)
			}
		})
	}
}

func TestUnFlattenMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected map[string]any
	}{
		{
			name: "simple unflatten",
			input: map[string]any{
				"user.id":   3,
				"user.name": "john",
			},
			expected: map[string]any{
				"user": map[string]any{
					"id":   3,
					"name": "john",
				},
			},
		},
		{
			name: "nested unflatten",
			input: map[string]any{
				"user.details.age":  30,
				"user.details.city": "Paris",
			},
			expected: map[string]any{
				"user": map[string]any{
					"details": map[string]any{
						"age":  30,
						"city": "Paris",
					},
				},
			},
		},
		{
			name: "flat map stays the same",
			input: map[string]any{
				"id":   1,
				"name": "Alice",
			},
			expected: map[string]any{
				"id":   1,
				"name": "Alice",
			},
		},
		{
			name: "deeply nested keys",
			input: map[string]any{
				"a.b.c.d": 42,
			},
			expected: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": map[string]any{
							"d": 42,
						},
					},
				},
			},
		},
		{
			name: "multiple root keys",
			input: map[string]any{
				"user.name":     "Bob",
				"user.age":      25,
				"location.city": "Berlin",
			},
			expected: map[string]any{
				"user": map[string]any{
					"name": "Bob",
					"age":  25,
				},
				"location": map[string]any{
					"city": "Berlin",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := UnFlattenMap(tc.input)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("UnFlattenMap() = %v, expected %v", got, tc.expected)
			}
		})
	}
}
