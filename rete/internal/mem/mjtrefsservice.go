package mem

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
	ri.nw = h.Nw
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
	nw     types.Network
}

func (ri *hdlTblIteratorImpl) HasNext() bool {
	return ri.curr != nil
}

func (ri *hdlTblIteratorImpl) Next() types.JoinTable {
	jtName := ri.curr.Value.(string)
	jT := ri.nw.GetJtService().GetJoinTable(jtName)
	ri.curr = ri.curr.Next()
	return jT
}

type RowIDIteratorImpl struct {
	jtName   string
	rowIdMap map[int]int
	kList    list.List
	curr     *list.Element
	nw       types.Network
}

func (ri *RowIDIteratorImpl) HasNext() bool {
	return ri.curr != nil
}

func (ri *RowIDIteratorImpl) Next() types.JoinTableRow {
	rowID := ri.curr.Value.(int)
	var jtRow types.JoinTableRow
	jT := ri.nw.GetJtService().GetJoinTable(ri.jtName)
	if jT != nil {
		jtRow = jT.GetRow(rowID)
	}
	ri.curr = ri.curr.Next()
	return jtRow
}

func (h *jtRefsServiceImpl) GetRowIterator(handle types.ReteHandle, jtName string) types.RowIterator {
	ri := RowIDIteratorImpl{}
	ri.jtName = jtName
	ri.kList = list.List{}
	ri.nw = h.Nw
	tblMap := h.tablesAndRows[handle.GetTupleKey().String()]
	if tblMap != nil {
		rowMap := tblMap[jtName]
		if rowMap != nil {
			ri.rowIdMap = rowMap
			for k, _ := range ri.rowIdMap {
				ri.kList.PushBack(k)
			}
		}
	}
	ri.curr = ri.kList.Front()
	return &ri
}
