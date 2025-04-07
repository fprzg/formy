package assert

import "testing"

func EqualError(t *testing.T, expected, actual error) {
}

func True(t testing.T, value bool, msgAndArgs ...interface{}) bool {
	return false
}
