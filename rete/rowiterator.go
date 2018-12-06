package rete

import "container/list"

type rowIterator interface {
	hasNext() bool
	next() joinTableRow
}

type rowIteratorImpl struct {
	table map[int]joinTableRow
	kList list.List
	curr  *list.Element
}

func (ri *rowIteratorImpl) hasNext() bool {
	return ri.curr != nil
}

func (ri *rowIteratorImpl) next() joinTableRow {
	id := ri.curr.Value.(int)
	val := ri.table[id]
	ri.curr = ri.curr.Next()
	return val
}
