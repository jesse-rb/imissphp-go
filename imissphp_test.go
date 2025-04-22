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
	type innerStruct struct {
		FieldA string `json:"field_a"`
		FieldB int    `json:"field_b"`
	}

	type testStruct struct {
		ID      int         `json:"id"`
		Name    string      `json:"name"`
		Details innerStruct `json:"details"`
		Tags    []string    `json:"tags"`
		Numbers [2]int      `json:"numbers"`
		private string      // unexported field
	}

	tests := []struct {
		name     string
		input    any
		expected map[string]any
	}{
		{
			name: "basic struct with nested struct and slice",
			input: testStruct{
				ID:   1,
				Name: "Alice",
				Details: innerStruct{
					FieldA: "abc",
					FieldB: 123,
				},
				Tags:    []string{"go", "dev"},
				Numbers: [2]int{7, 8},
				private: "should be ignored",
			},
			expected: map[string]any{
				"id":   1,
				"name": "Alice",
				"details": map[string]any{
					"field_a": "abc",
					"field_b": 123,
				},
				"tags": map[string]any{
					"0": "go",
					"1": "dev",
				},
				"numbers": map[string]any{
					"0": 7,
					"1": 8,
				},
			},
		},
		{
			name:  "plain slice input",
			input: []int{10, 20, 30},
			expected: map[string]any{
				"0": 10,
				"1": 20,
				"2": 30,
			},
		},
		{
			name:  "plain array input",
			input: [3]string{"a", "b", "c"},
			expected: map[string]any{
				"0": "a",
				"1": "b",
				"2": "c",
			},
		},
		{
			name: "simple map[string]any",
			input: map[string]any{
				"foo": 42,
				"bar": "baz",
			},
			expected: map[string]any{
				"foo": 42,
				"bar": "baz",
			},
		},
		{
			name:     "non-convertible value (int)",
			input:    123,
			expected: map[string]any{},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: map[string]any{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ToMap(tc.input)

			// Compare JSON since the function's use of JSON.Marshal not preserve exact types
			expectedJSON, err := json.Marshal(tc.expected)
			if err != nil {
				t.Fatalf("error marshaling expected: %v", err)
			}

			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("error marshaling result: %v", err)
			}

			if string(expectedJSON) != string(gotJSON) {
				t.Errorf("Test %s failed\nExpected JSON: %s\nGot JSON: %s", tc.name, string(expectedJSON), string(gotJSON))
			}
		})
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
