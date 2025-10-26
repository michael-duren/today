package main

import "testing"

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("execpted nil error, got %v", err)
	}
}

func AssertTrue(t testing.TB, actual bool) {
	t.Helper()

	if !actual {
		t.Errorf("Expected val %v was false, expected true", actual)
	}
}

func AssertFalse(t testing.TB, actual bool) {
	t.Helper()

	if actual {
		t.Errorf("Expected val %v was true, expected valse", actual)
	}
}

func AssertEqual(t testing.TB, expected, actual any) {
	t.Helper()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

type Stringer interface {
	String() string
}

func AseertEqualString(t testing.TB, expected, actual Stringer) {
	t.Helper()

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func AssertError(t testing.TB, expected, actual error) {
	t.Helper()
	if expected == nil {
		t.Error("wanted an error but didn't get one")
	}
	if expected != actual {
		t.Errorf("expected error %s, got %s", expected, actual)
	}
}
