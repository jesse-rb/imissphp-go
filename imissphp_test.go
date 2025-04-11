package imissphpgo

import "testing"

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
