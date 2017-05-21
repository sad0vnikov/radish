package db

import (
	"errors"
	"fmt"
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

func TestGettingChildrenFromKeys(t *testing.T) {
	keys := []string{"a", "b", "c", "d", "e"}
	keysMask := "*"
	delimiter := ":"

	result := getChildrenFromKeys(keys, keysMask, delimiter)

	expectedResult := []KeyTreeNode{
		KeyTreeNode{Name: "a", HasChildren: false},
		KeyTreeNode{Name: "b", HasChildren: false},
		KeyTreeNode{Name: "c", HasChildren: false},
		KeyTreeNode{Name: "d", HasChildren: false},
		KeyTreeNode{Name: "e", HasChildren: false},
	}

	if err := compareTreeNodeSlices(result, expectedResult); err != nil {
		t.Error(err)
	}

	keys = []string{"a:c", "b", "c", "d:r:f", "e"}
	keysMask = "*"
	delimiter = ":"

	result = getChildrenFromKeys(keys, keysMask, delimiter)

	expectedResult = []KeyTreeNode{
		KeyTreeNode{Name: "a", HasChildren: true},
		KeyTreeNode{Name: "b", HasChildren: false},
		KeyTreeNode{Name: "c", HasChildren: false},
		KeyTreeNode{Name: "d", HasChildren: true},
		KeyTreeNode{Name: "e", HasChildren: false},
	}

	if err := compareTreeNodeSlices(result, expectedResult); err != nil {
		t.Error(err)
	}

	keys = []string{"d:r:f", "d:r:e"}
	keysMask = "d:*"
	delimiter = ":"

	result = getChildrenFromKeys(keys, keysMask, delimiter)

	expectedResult = []KeyTreeNode{
		KeyTreeNode{Name: "r", HasChildren: true},
	}

	if err := compareTreeNodeSlices(result, expectedResult); err != nil {
		t.Error(err)
	}

	keys = []string{"d:r:f", "d:r:e"}
	keysMask = "d:r:*"
	delimiter = ":"

	result = getChildrenFromKeys(keys, keysMask, delimiter)

	expectedResult = []KeyTreeNode{
		KeyTreeNode{Name: "f", HasChildren: false},
		KeyTreeNode{Name: "e", HasChildren: false},
	}

	if err := compareTreeNodeSlices(result, expectedResult); err != nil {
		t.Error(err)
	}

}

func compareTreeNodeSlices(a, b []KeyTreeNode) error {
	if a == nil && b == nil {
		return nil
	}

	if a == nil && b != nil {
		return errors.New("first slice is nil, but second is not")
	}

	if a != nil && b == nil {
		return errors.New("second slice is nil, but first is not")
	}

	for len(a) != len(b) {
		return fmt.Errorf("the slices %v and %v are different", a, b)
	}

	for i := range a {
		if a[i].HasChildren != b[i].HasChildren || a[i].Name != b[i].Name {
			return errors.New("elements #" + string(i) + " are different in lists")
		}
	}

	return nil
}
