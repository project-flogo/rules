package mem

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
)

type joinTableImpl struct {
	types.NwElemIdImpl
	table map[int]types.JoinTableRow
	idr   []model.TupleType
	rule  model.Rule
	name string
}

func newJoinTableImpl (nw types.Network, rule model.Rule, identifiers []model.TupleType, name string) types.JoinTable {
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

	jt.table[row.GetID()] = row
	for i := 0; i < len(row.GetHandles()); i++ {
		handle := row.GetHandles()[i]
		jt.Nw.GetJtRefService().AddEntry(handle, jt.name, row.GetID())
	}
	return row
}

func (jt *joinTableImpl) RemoveRow(rowID int) types.JoinTableRow {
	row, found := jt.table[rowID]
	if found {
		delete(jt.table, rowID)
		return row
	}
	return nil
}

func (jt *joinTableImpl) GetRowCount() int {
	return len(jt.table)
}

func (jt *joinTableImpl) GetRule() model.Rule {
	return jt.rule
}

func (jt *joinTableImpl) GetRowIterator() types.RowIterator {
	return newRowIterator(jt.table)
}

func (jt *joinTableImpl) GetRow(rowID int) types.JoinTableRow {
	return jt.table[rowID]
}

func (jt *joinTableImpl) GetName() string {
	return jt.name
}