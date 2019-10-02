package mem

import (
	"container/list"
	"github.com/project-flogo/rules/rete/internal/types"
)

type jtRefsServiceImpl struct {
	//keys are jointable-ids and values are lists of row-ids in the corresponding join table
	types.NwServiceImpl
	tablesAndRows map[string]map[int]string
}

func NewJoinTableRefsInHdlImpl(nw types.Network, config map[string]interface{}) types.JtRefsService {
	hdlJt := jtRefsServiceImpl{}
	hdlJt.Nw = nw
	hdlJt.tablesAndRows = make(map[string]map[int]string)
	return &hdlJt
}

func (h *jtRefsServiceImpl) Init() {

}

func (h *jtRefsServiceImpl) AddEntry(handle types.ReteHandle, jtName string, rowID int) {
	tblMap, found := h.tablesAndRows[handle.GetTupleKey().String()]
	if !found {
		tblMap = make(map[int]string)
		h.tablesAndRows[handle.GetTupleKey().String()] = tblMap
	}
	tblMap[rowID] = jtName
}

func (h *jtRefsServiceImpl) RemoveEntry(handle types.ReteHandle, jtName string, rowID int) {
	tblMap, found := h.tablesAndRows[handle.GetTupleKey().String()]
	if found {
		delete(tblMap, rowID)
	}
}

type hdlRefsRowIterator struct {
	rowIdMap  map[int]string
	kList     list.List
	curr      *list.Element
	currRowId int
	nw        types.Network
}

func (ri *hdlRefsRowIterator) HasNext() bool {
	return ri.curr != nil
}

func (ri *hdlRefsRowIterator) Next() (types.JoinTableRow, types.JoinTable) {
	rowID := ri.curr.Value.(int)
	ri.currRowId = rowID
	ri.curr = ri.curr.Next()
	jT := ri.nw.GetJtService().GetJoinTable(ri.rowIdMap[rowID])
	if jT != nil {
		return jT.GetRow(rowID), jT
	}
	return nil, jT
}

func (ri *hdlRefsRowIterator) Remove() {
	delete(ri.rowIdMap, ri.currRowId)
}

func (h *jtRefsServiceImpl) GetRowIterator(handle types.ReteHandle) types.JointableIterator {
	ri := hdlRefsRowIterator{}
	ri.kList = list.List{}
	ri.nw = h.Nw
	tblMap := h.tablesAndRows[handle.GetTupleKey().String()]
	ri.rowIdMap = tblMap
	for rowID := range tblMap {
		ri.kList.PushBack(rowID)
	}
	ri.curr = ri.kList.Front()
	return &ri
}
