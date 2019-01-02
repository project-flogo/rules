package redis

import (
	"container/list"
	"github.com/project-flogo/rules/rete/internal/types"
)

type jtRefsServiceImpl struct {
	//keys are jointable-ids and values are lists of row-ids in the corresponding join table
	types.NwServiceImpl

	tablesAndRows map[string]map[string]map[int]int
}

func NewJoinTableRefsInHdlImpl(nw types.Network, config map[string]interface{}) types.JtRefsService {
	hdlJt := jtRefsServiceImpl{}
	hdlJt.Nw = nw
	hdlJt.tablesAndRows = make(map[string]map[string]map[int]int)
	return &hdlJt
}

func (h *jtRefsServiceImpl) Init() {

}

func (h *jtRefsServiceImpl) AddEntry(handle types.ReteHandle, jtName string, rowID int) {

	tblMap, found := h.tablesAndRows[handle.GetTupleKey().String()]

	if !found {
		tblMap = make(map[string]map[int]int)
		h.tablesAndRows[handle.GetTupleKey().String()] = tblMap
	}

	rowsForJoinTable, found := tblMap[jtName]
	if !found {
		rowsForJoinTable = make(map[int]int)
		tblMap[jtName] = rowsForJoinTable
	}
	rowsForJoinTable[rowID] = rowID
}

func (h *jtRefsServiceImpl) RemoveEntry(handle types.ReteHandle, jtName string) {
	tblMap, found := h.tablesAndRows[handle.GetTupleKey().String()]
	if found {
		delete(tblMap, jtName)
	}
}

func (h *jtRefsServiceImpl) RemoveRowEntry(handle types.ReteHandle, jtName string, rowID int) {
	tblMap, found := h.tablesAndRows[handle.GetTupleKey().String()]
	if found {
		rowIDs, fnd := tblMap[jtName]
		if fnd {
			delete(rowIDs, rowID)
		}
	}
}

func (h *jtRefsServiceImpl) GetIterator(handle types.ReteHandle) types.HdlTblIterator {
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
	tblMap map[string]map[int]int
	kList  list.List
	curr   *list.Element
}

func (ri *hdlTblIteratorImpl) HasNext() bool {
	return ri.curr != nil
}

func (ri *hdlTblIteratorImpl) Next() (string, map[int]int) {
	id := ri.curr.Value.(string)
	lst := ri.tblMap[id]
	ri.curr = ri.curr.Next()
	return id, lst
}
