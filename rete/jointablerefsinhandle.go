package rete

import "container/list"

type joinTableRefsInHdl interface {
	addEntry(jointTableID int, rowID int)
	removeEntry(jointTableID int)
}

type joinTableRefsInHdlImpl struct {
	//keys are jointable-ids and values are lists of row-ids in the corresponding join table
	tablesAndRows map[int]*list.List
}

func newJoinTableRefsInHdlImpl() joinTableRefsInHdl {
	hdlJt := joinTableRefsInHdlImpl{}
	hdlJt.tablesAndRows = make(map[int]*list.List)
	return &hdlJt
}

func (h *joinTableRefsInHdlImpl) addEntry(jointTableID int, rowID int) {
	rowsForJoinTable := h.tablesAndRows[jointTableID]
	if rowsForJoinTable == nil {
		rowsForJoinTable = list.New()
		h.tablesAndRows[jointTableID] = rowsForJoinTable
	}
	rowsForJoinTable.PushBack(rowID)
}

func (h *joinTableRefsInHdlImpl) removeEntry(jointTableID int) {
	delete(h.tablesAndRows, jointTableID)
}

type hdlTblIterator interface {
	hasNext() bool
	next() (int, *list.List)
}

type hdlTblIteratorImpl struct {
	hdlJtImpl *joinTableRefsInHdlImpl
	kList     list.List
	curr      *list.Element
}

func (ri *hdlTblIteratorImpl) hasNext() bool {
	return ri.curr != nil
}

func (ri *hdlTblIteratorImpl) next() (int, *list.List) {
	id := ri.curr.Value.(int)
	lst := ri.hdlJtImpl.tablesAndRows[id]
	ri.curr = ri.curr.Next()
	return id, lst
}

func (hdl *reteHandleImpl) newHdlTblIterator() hdlTblIterator {
	ri := hdlTblIteratorImpl{}
	ri.hdlJtImpl = hdl.jtRefs.(*joinTableRefsInHdlImpl)
	ri.kList = list.List{}
	for k, _ := range ri.hdlJtImpl.tablesAndRows {
		ri.kList.PushBack(k)
	}
	ri.curr = ri.kList.Front()
	return &ri
}
