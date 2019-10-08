package mem

import (
	"container/list"
	"context"
	"sync"

	"github.com/project-flogo/rules/rete/internal/types"
)

type jtRowsImpl struct {
	rows map[int]string
	sync.RWMutex
}

type jtRefsServiceImpl struct {
	types.NwServiceImpl
	tablesAndRows map[string]*jtRowsImpl
	sync.RWMutex
}

func NewJoinTableRefsInHdlImpl(nw types.Network, config map[string]interface{}) types.JtRefsService {
	hdlJt := jtRefsServiceImpl{}
	hdlJt.Nw = nw
	hdlJt.tablesAndRows = make(map[string]*jtRowsImpl)
	return &hdlJt
}

func (h *jtRefsServiceImpl) Init() {

}

func (h *jtRefsServiceImpl) AddEntry(handle types.ReteHandle, jtName string, rowID int) {
	key := handle.GetTupleKey().String()
	h.Lock()
	defer h.Unlock()
	tblMap, found := h.tablesAndRows[key]
	if !found {
		tblMap = &jtRowsImpl{
			rows: make(map[int]string),
		}
		h.tablesAndRows[handle.GetTupleKey().String()] = tblMap
	}
	tblMap.Lock()
	defer tblMap.Unlock()
	tblMap.rows[rowID] = jtName
}

func (h *jtRefsServiceImpl) RemoveEntry(handle types.ReteHandle, jtName string, rowID int) {
	key := handle.GetTupleKey().String()
	h.Lock()
	defer h.Unlock()
	tblMap, found := h.tablesAndRows[key]
	if found {
		tblMap.Lock()
		defer tblMap.Unlock()
		delete(tblMap.rows, rowID)
		if len(tblMap.rows) == 0 {
			delete(h.tablesAndRows, key)
		}
	}
}

type hdlRefsRowIterator struct {
	ctx     context.Context
	refs    *jtRefsServiceImpl
	key     string
	rows    *jtRowsImpl
	list    list.List
	current *list.Element
	rowID   int
	nw      types.Network
}

func (ri *hdlRefsRowIterator) HasNext() bool {
	return ri.current != nil
}

func (ri *hdlRefsRowIterator) Next() (types.JoinTableRow, types.JoinTable) {
	rowID := ri.current.Value.(int)
	ri.rowID = rowID
	ri.current = ri.current.Next()
	ri.rows.RLock()
	defer ri.rows.RUnlock()
	jT := ri.nw.GetJtService().GetJoinTable(ri.rows.rows[rowID])
	if jT != nil {
		return jT.GetRow(ri.ctx, rowID), jT
	}
	return nil, jT
}

func (ri *hdlRefsRowIterator) Remove() {
	ri.rows.Lock()
	defer ri.rows.Unlock()
	delete(ri.rows.rows, ri.rowID)
	if len(ri.rows.rows) == 0 {
		ri.refs.Lock()
		defer ri.refs.Unlock()
		delete(ri.refs.tablesAndRows, ri.key)
	}
}

func (h *jtRefsServiceImpl) GetRowIterator(ctx context.Context, handle types.ReteHandle) types.JointableIterator {
	key := handle.GetTupleKey().String()
	ri := hdlRefsRowIterator{
		ctx:  ctx,
		refs: h,
		key:  key,
		nw:   h.Nw,
	}
	h.RLock()
	defer h.RUnlock()
	tblMap := h.tablesAndRows[key]
	if tblMap == nil {
		return &ri
	}
	ri.rows = tblMap
	for rowID := range tblMap.rows {
		ri.list.PushBack(rowID)
	}
	ri.current = ri.list.Front()
	return &ri
}
