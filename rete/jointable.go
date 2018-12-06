package rete

import (
	"github.com/project-flogo/rules/common/model"
	"container/list"
)

type joinTable interface {
	addRow(row joinTableRow) //list of Tuples
	getID() int
	len() int
	//getMap() map[joinTableRow]joinTableRow
	removeRow(rowID int)
	getRule() model.Rule
	iterator() rowIterator
}

type joinTableImpl struct {
	id    int
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
	jt.id = nw.incrementAndGetId()
	jt.idr = identifiers
	jt.table = map[int]joinTableRow{}
	jt.rule = rule
}

func (jt *joinTableImpl) getID() int {
	return jt.id
}

func (jt *joinTableImpl) addRow(row joinTableRow) {
	jt.table[row.getID()] = row
	for i := 0; i < len(row.getHandles()); i++ {
		handle := row.getHandles()[i]
		handle.addJoinTableRowRef(row, jt)
	}
}

func (jt *joinTableImpl) removeRow(rowID int) {
	delete(jt.table, rowID)
}

func (jt *joinTableImpl) len() int {
	return len(jt.table)
}

//func (jt *joinTableImpl) getMap() map[joinTableRow]joinTableRow {
//	return jt.table
//}

func (jt *joinTableImpl) getRule() model.Rule {
	return jt.rule
}

func (jt *joinTableImpl) iterator() rowIterator {
	ri := rowIteratorImpl{}
	ri.table = jt.table
	ri.kList = list.List{}
	for k, _:= range jt.table {
		ri.kList.PushBack(k)
	}
	ri.curr = ri.kList.Front()
	return &ri
}