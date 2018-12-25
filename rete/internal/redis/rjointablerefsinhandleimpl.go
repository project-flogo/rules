package redis

import (
	"container/list"
	"github.com/project-flogo/rules/rete/internal/types"
)

type joinTableRefsInHdlImpl struct {
	//keys are jointable-ids and values are lists of row-ids in the corresponding join table
	tablesAndRows map[string]map[string]*list.List
}

func NewJoinTableRefsInHdlImpl(config map[string]interface{}) types.JtRefsService {
	hdlJt := joinTableRefsInHdlImpl{}
	hdlJt.tablesAndRows = make(map[string]map[string]*list.List)
	return &hdlJt
}

func (h *joinTableRefsInHdlImpl) Init() {

}

func (h *joinTableRefsInHdlImpl) AddEntry(handle types.ReteHandle, jtName string, rowID int) {

	tblMap, found := h.tablesAndRows[handle.GetTupleKey().String()]

	if !found {
		tblMap = make(map[string]*list.List)
		h.tablesAndRows[handle.GetTupleKey().String()] = tblMap
	}

	rowsForJoinTable, found := tblMap[jtName]
	if !found {
		rowsForJoinTable = list.New()
		tblMap[jtName] = rowsForJoinTable
	}
	rowsForJoinTable.PushBack(rowID)
}

func (h *joinTableRefsInHdlImpl) RemoveEntry(handle types.ReteHandle, jtName string) {
	tblMap, found := h.tablesAndRows[handle.GetTupleKey().String()]
	if found {
		delete(tblMap, jtName)
	}
}

func (h *joinTableRefsInHdlImpl) RemoveRowEntry(handle types.ReteHandle, jtName string, rowID int) {
	tblMap, found := h.tablesAndRows[handle.GetTupleKey().String()]
	if found {
		rowIDs, fnd := tblMap[jtName]
		if fnd {
			for e:= rowIDs.Front(); e != nil; e = e.Next() {
				rowIDInList := e.Value.(int)
				if rowID == rowIDInList {
					rowIDs.Remove(e)
					return
				}
			}
		}
	}
}

func (h *joinTableRefsInHdlImpl) GetIterator(handle types.ReteHandle) types.HdlTblIterator {
	ri := hdlTblIteratorImpl{}
	//ri.hdlJtImpl = h
	ri.kList = list.List{}

	tblMap, found := h.tablesAndRows[handle.GetTupleKey().String()]
	if found {
		ri.tblMap = tblMap
		for k, _ := range tblMap {
			ri.kList.PushBack(k)
		}
	}
	ri.curr = ri.kList.Front()
	return &ri
}

type hdlTblIteratorImpl struct {
	tblMap map[string]*list.List
	kList  list.List
	curr   *list.Element
}

func (ri *hdlTblIteratorImpl) HasNext() bool {
	return ri.curr != nil
}

func (ri *hdlTblIteratorImpl) Next() (string, *list.List) {
	id := ri.curr.Value.(string)
	lst := ri.tblMap[id]
	ri.curr = ri.curr.Next()
	return id, lst
}
