package redis

import (
	"container/list"
	"github.com/project-flogo/rules/rete/internal/types"
)

type joinTableRefsInHdlImpl struct {
	//keys are jointable-ids and values are lists of row-ids in the corresponding join table
	tablesAndRows map[string]*list.List
}

func NewJoinTableRefsInHdlImpl(config map[string]interface{}) types.JtRefsService {
	hdlJt := joinTableRefsInHdlImpl{}
	hdlJt.tablesAndRows = make(map[string]*list.List)
	return &hdlJt
}

func (h *joinTableRefsInHdlImpl) Init() {

}

func (h *joinTableRefsInHdlImpl) AddEntry(handle types.ReteHandle, jtName string, rowID int) {
	rowsForJoinTable := h.tablesAndRows[jtName]
	if rowsForJoinTable == nil {
		rowsForJoinTable = list.New()
		h.tablesAndRows[jtName] = rowsForJoinTable
	}
	rowsForJoinTable.PushBack(rowID)
}

func (h *joinTableRefsInHdlImpl) RemoveEntry(handle types.ReteHandle, jtName string) {
	delete(h.tablesAndRows, jtName)
}

func (h *joinTableRefsInHdlImpl) GetIterator(handle types.ReteHandle) types.HdlTblIterator {
	ri := hdlTblIteratorImpl{}
	ri.hdlJtImpl = h
	ri.kList = list.List{}
	for k, _ := range ri.hdlJtImpl.tablesAndRows {
		ri.kList.PushBack(k)
	}
	ri.curr = ri.kList.Front()
	return &ri
}

type hdlTblIteratorImpl struct {
	hdlJtImpl *joinTableRefsInHdlImpl
	kList     list.List
	curr      *list.Element
}

func (ri *hdlTblIteratorImpl) HasNext() bool {
	return ri.curr != nil
}

func (ri *hdlTblIteratorImpl) Next() (string, *list.List) {
	id := ri.curr.Value.(string)
	lst := ri.hdlJtImpl.tablesAndRows[id]
	ri.curr = ri.curr.Next()
	return id, lst
}
