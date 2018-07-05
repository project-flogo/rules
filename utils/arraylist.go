package utils

//ArrayList A simple array list interface using slices
type ArrayList interface {
	Add(interface{})
	Remove(interface{}) bool
	Get(int) interface{}
	RemoveAt(int) interface{}
	Len() int
	Contains(interface{}) bool
	ForEach(fn ArrayListEntryHandler, context []interface{})
}

//Callback to iterate on list entries
type ArrayListEntryHandler func(entry interface{}, context []interface{})

type arrayList struct {
	items []interface{}
}

//NewArrayList Allocate a new arrayList
func NewArrayList() ArrayList {
	arrayListVar := arrayList{}
	return &arrayListVar
}

func (arrayListVar *arrayList) Add(item interface{}) {
	arrayListVar.items = append(arrayListVar.items, item)
}

func (arrayListVar *arrayList) Remove(item interface{}) bool {
	newslice := []interface{}{}
	removed := false
	i := 0
	for ; i < arrayListVar.Len(); i++ {
		current := arrayListVar.Get(i)
		if current != nil && current == item {
			removed = true
			break
		}
	}
	if removed {
		newslice = append(newslice, arrayListVar.items[:i]...)

		if arrayListVar.Len() > i+1 {
			newslice = append(newslice, arrayListVar.items[i+1:]...)
		}
		arrayListVar.items = newslice
	}
	return removed
}

func (arrayListVar *arrayList) Get(index int) interface{} {
	if index < 0 || index >= len(arrayListVar.items) {
		return nil
	}
	return arrayListVar.items[index]
}

func (arrayListVar *arrayList) RemoveAt(index int) interface{} {
	newslice := []interface{}{}
	if index < 0 || index >= arrayListVar.Len() {
		return nil
	}

	newslice = append(newslice, arrayListVar.items[:index]...)
	if arrayListVar.Len() > index+1 {
		newslice = append(newslice, arrayListVar.items[index+1:]...)
	}
	removed := arrayListVar.items[index]
	arrayListVar.items = newslice
	return removed
}

func (arrayListVar *arrayList) Len() int {
	return len(arrayListVar.items)
}

func (arrayListVar *arrayList) Contains(item interface{}) bool {
	if item == nil {
		return false
	}
	for i := 0; i < arrayListVar.Len(); i++ {
		if item == arrayListVar.Get(i) {
			return true
		}
	}
	return false

}

func (arrayListVar *arrayList) ForEach(entryHandler ArrayListEntryHandler, context []interface{}) {
	for i := 0; i < arrayListVar.Len(); i++ {
		entryHandler(arrayListVar.Get(i), context)
	}
}
