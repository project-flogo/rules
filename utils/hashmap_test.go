package utils

import (
	"fmt"
	"testing"
)

func TestHashMap(t *testing.T) {

	m := NewLinkedHashMap()

	m.Put("a", "a")
	m.Put("b", "b")
	m.Put("c", "c")

	m.ForEach(HandleEntry, nil)

	m.Remove("b")
	m.ForEach(HandleEntry, nil)
}

func HandleEntry(key string, value interface{}, context []interface{}) {
	myVal := value.(string)
	fmt.Println("key:" + key + ", value:" + myVal)
}
