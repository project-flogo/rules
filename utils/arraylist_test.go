package utils

import (
	"fmt"
	"strconv"
	"testing"
)

func TestArrayList(t *testing.T) {

	myList := NewArrayList()

	myList.Add(1)
	myList.Add(2)
	myList.Add(3)
	myList.Add(4)

	myList.Remove(4)

	myList.RemoveAt(7)
	myList.RemoveAt(0)
	myList.RemoveAt(6)

	var x = make([]interface{}, 1)
	x[0] = "hi"

	myList.ForEach(PrintItem, x)

}

func PrintItem(entry interface{}, context []interface{}) {
	i := entry.(int)
	str1 := len(context)
	x := context[0].(string)

	str := strconv.Itoa(i) + ", " + strconv.Itoa(str1)

	fmt.Println(str + ", " + x)
}
