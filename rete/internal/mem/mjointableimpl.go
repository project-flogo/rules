package mem

import (
	"sync"

	"container/list"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
)

type joinTableImpl struct {
	types.NwElemIdImpl
	table map[int]types.JoinTableRow
	idr   []model.TupleType
	rule  model.Rule
	name  string
	sync.RWMutex
}

func newJoinTableImpl(nw types.Network, rule model.Rule, identifiers []model.TupleType, name string) types.JoinTable {
	jt := joinTableImpl{}
	jt.initJoinTableImpl(nw, rule, identifiers, name)
	return &jt
}

func (jt *joinTableImpl) initJoinTableImpl(nw types.Network, rule model.Rule, identifiers []model.TupleType, name string) {
	jt.SetID(nw)
	jt.idr = identifiers
	jt.table = map[int]types.JoinTableRow{}
	jt.rule = rule
	jt.name = name
}

func (jt *joinTableImpl) AddRow(handles []types.ReteHandle) types.JoinTableRow {
	row := newJoinTableRow(handles, jt.Nw)
	for i := 0; i < len(row.GetHandles()); i++ {
		handle := row.GetHandles()[i]
		jt.Nw.GetJtRefService().AddEntry(handle, jt.name, row.GetID())
	}
	jt.Lock()
	defer jt.Unlock()
	jt.table[row.GetID()] = row
	return row
}

func (jt *joinTableImpl) RemoveRow(rowID int) types.JoinTableRow {
	jt.Lock()
	defer jt.Unlock()
	row, found := jt.table[rowID]
	if found {
		delete(jt.table, rowID)
		return row
	}
	return nil
}

func (jt *joinTableImpl) GetRowCount() int {
	jt.RLock()
	defer jt.RUnlock()
	return len(jt.table)
}

func (jt *joinTableImpl) GetRule() model.Rule {
	return jt.rule
}

func (jt *joinTableImpl) GetRowIterator() types.JointableRowIterator {
	return newRowIterator(jt)
}

func (jt *joinTableImpl) GetRow(rowID int) types.JoinTableRow {
	jt.RLock()
	defer jt.RUnlock()
	return jt.table[rowID]
}

func (jt *joinTableImpl) GetName() string {
	return jt.name
}

func (jt *joinTableImpl) RemoveAllRows() {
	rowIter := jt.GetRowIterator()
	for rowIter.HasNext() {
		row := rowIter.Next()
		//first, from jTable, remove row
		jt.RemoveRow(row.GetID())
		for _, hdl := range row.GetHandles() {
			jt.Nw.GetJtRefService().RemoveEntry(hdl, jt.GetName(), row.GetID())
		}
		//Delete the rowRef itself
		rowIter.Remove()
	}
}

type rowIteratorImpl struct {
	table   *joinTableImpl
	list    list.List
	currKey int
	curr    *list.Element
}

func newRowIterator(jt *joinTableImpl) types.JointableRowIterator {
	ri := rowIteratorImpl{
		table: jt,
	}
	jt.RLock()
	defer jt.RUnlock()
	for k := range jt.table {
		ri.list.PushBack(k)
	}
	ri.curr = ri.list.Front()
	return &ri
}

func (ri *rowIteratorImpl) HasNext() bool {
	return ri.curr != nil
}

func (ri *rowIteratorImpl) Next() types.JoinTableRow {
	ri.currKey = ri.curr.Value.(int)
	ri.curr = ri.curr.Next()
	ri.table.RLock()
	defer ri.table.RUnlock()
	return ri.table.table[ri.currKey]
}

func (ri *rowIteratorImpl) Remove() {
	ri.table.Lock()
	defer ri.table.Unlock()
	delete(ri.table.table, ri.currKey)
}
