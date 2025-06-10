package imissphp

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"
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
			expected: map[string]any{
				"user.details.age":  30,
				"user.details.city": "Paris",
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
			expected: map[string]any{
				"a.b.c.d": 42,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := FlattenMap(tc.input)
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

func TestMapValues(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []int
	}{
		{
			input: map[string]int{
				"one": 1, "awljdnjwnd": 983, "-3": -33,
			},
			expected: []int{
				1, 983, -33,
			},
		},
	}

	for _, test := range tests {
		actual := MapValues(test.input)
		if len(actual) != len(test.expected) {
			t.Fatalf("Expected actual length: %v to match expected length: %v", len(actual), len(test.expected))
		}

		slices.Sort(actual)
		slices.Sort(test.expected)

		for i := 0; i < len(actual); i++ {
			if actual[i] != test.expected[i] {
				t.Fatalf("Expected actual[%v]: %v to equal expected[%v]: %v", i, actual[i], i, test.expected[i])
			}
		}
	}
}

func TestMapKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []string
	}{
		{
			input: map[string]int{
				"abc": 1, "BC-2": 32, "23": 9982,
			},
			expected: []string{
				"abc", "BC-2", "23",
			},
		},
	}

	for _, test := range tests {
		actual := MapKeys(test.input)
		if len(actual) != len(test.expected) {
			t.Fatalf("Expected actual length: %v to match expected length: %v", len(actual), len(test.expected))
		}

		slices.Sort(actual)
		slices.Sort(test.expected)

		for i := 0; i < len(actual); i++ {
			if actual[i] != test.expected[i] {
				t.Fatalf("Expected actual[%v]: %v to equal expected[%v]: %v", i, actual[i], i, test.expected[i])
			}
		}
	}
}
