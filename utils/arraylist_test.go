package utils

import (
	"container/list"
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

func TestOne(t *testing.T) {

	l := list.New()

	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(4)
	l.PushBack(5)

	for e := l.Front(); e != nil; {
		i := e.Value.(int)
		y := e.Next()

		if i%2 != 0 {
			l.Remove(e)
		} else {
			fmt.Printf("%d\n", i)
		}
		e = y
	}

	for e := l.Front(); e != nil; e = e.Next() {
		i := e.Value.(int)
		fmt.Printf("%d\n", i)
	}
}
