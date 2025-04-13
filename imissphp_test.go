package imissphpgo

import (
	"log"
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
