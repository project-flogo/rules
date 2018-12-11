package rete

import "github.com/project-flogo/rules/common/model"

type joinTableImpl struct {
	nwElemIdImpl
	table map[int]joinTableRow
	idr   []model.TupleType
	rule  model.Rule
}

func newJoinTable(nw Network, rule model.Rule, identifiers []model.TupleType) joinTable {
	jT := joinTableImpl{}
	jT.initJoinTableImpl(nw, rule, identifiers)

	//add it to all join tables collection before returning
	reteNw := nw.(*reteNetworkImpl)
	reteNw.allJoinTables[jT.getID()] = &jT
	return &jT
}

func (jt *joinTableImpl) initJoinTableImpl(nw Network, rule model.Rule, identifiers []model.TupleType) {
	jt.setID(nw)
	jt.idr = identifiers
	jt.table = map[int]joinTableRow{}
	jt.rule = rule
}

func (jt *joinTableImpl) addRow(handles []reteHandle) joinTableRow {

	row := newJoinTableRow(handles, jt.nw)

	jt.table[row.getID()] = row
	for i := 0; i < len(row.getHandles()); i++ {
		handle := row.getHandles()[i]
		handle.addJoinTableRowRef(row, jt)
	}
	return row
}

func (jt *joinTableImpl) removeRow(rowID int) joinTableRow {
	row, found := jt.table[rowID]
	if found {
		delete(jt.table, rowID)
		return row
	}
	return nil
}

func (jt *joinTableImpl) getRowCount() int {
	return len(jt.table)
}

func (jt *joinTableImpl) getRule() model.Rule {
	return jt.rule
}

func (jt *joinTableImpl) getRowIterator() rowIterator {
	return newRowIterator(jt.table)
}

func (jt *joinTableImpl) getRow(rowID int) joinTableRow {
	return jt.table[rowID]
}
