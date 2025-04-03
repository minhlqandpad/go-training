package utils

import (
	"reflect"
	"testing"
)

func AssertCorrectResult[T any](t testing.TB, result, expected T) {
	t.Helper()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected '%v' but got result '%v'", expected, result)
	}
}

func AssertError(t testing.TB, err error, shouldError bool) {
	t.Helper()
	if shouldError && err == nil {
		t.Error("Expected to get an error, but got no error at all")
	} else if !shouldError && err != nil {
		t.Errorf("Expected to run without error, but got error \"%v\" instead", err)
	}
}
