package db

import (
	"testing"

	"github.com/rafaeljusto/redigomock"
)

func TestGetKeysByMask(t *testing.T) {

	conn := redigomock.NewConn()
	connector = &MockedConnections{ConnectionMock: conn}

	expectedResult := []interface{}{}

	expectedResult = append(expectedResult, interface{}([]byte("aa")))
	expectedResult = append(expectedResult, interface{}([]byte("ab")))
	expectedResult = append(expectedResult, interface{}([]byte("ac")))

	conn.Command("KEYS", "a*").
		Expect(expectedResult)

	result, err := FindKeysByMask("server1", "a*")
	if err != nil {
		t.Error(err)
	}

	if !checkSlicesAreEqual(result, []string{"aa", "ab", "ac"}) {
		t.Errorf("got invalid keys list: %v, expected %v", result, expectedResult)
	}
}

func checkSlicesAreEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
