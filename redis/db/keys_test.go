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

	result, err := FindKeysByMask("server1", 0, "a*")
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
		KeyTreeNode{Name: "a", Key: "a", HasChildren: false},
		KeyTreeNode{Name: "b", Key: "b", HasChildren: false},
		KeyTreeNode{Name: "c", Key: "c", HasChildren: false},
		KeyTreeNode{Name: "d", Key: "d", HasChildren: false},
		KeyTreeNode{Name: "e", Key: "e", HasChildren: false},
	}

	if err := compareTreeNodeSlices(result, expectedResult); err != nil {
		t.Error(err)
	}

	keys = []string{"a:c", "b", "c", "d:r:f", "e"}
	keysMask = "*"
	delimiter = ":"

	result = getChildrenFromKeys(keys, keysMask, delimiter)

	expectedResult = []KeyTreeNode{
		KeyTreeNode{Name: "a", Key: "", HasChildren: true},
		KeyTreeNode{Name: "b", Key: "b", HasChildren: false},
		KeyTreeNode{Name: "c", Key: "c", HasChildren: false},
		KeyTreeNode{Name: "d", Key: "", HasChildren: true},
		KeyTreeNode{Name: "e", Key: "e", HasChildren: false},
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
		KeyTreeNode{Name: "f", Key: "f", HasChildren: false},
		KeyTreeNode{Name: "e", Key: "e", HasChildren: false},
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
		return fmt.Errorf("the slices %v and %v have different length", a, b)
	}

	for i := range a {
		prs := false
		for j := range b {
			if a[i].HasChildren == b[j].HasChildren && a[i].Name == b[j].Name {
				prs = true
				break
			}
		}
		if !prs {
			return fmt.Errorf("element: %#v not found in list %#v", a[i], b)
		}
	}

	return nil
}
