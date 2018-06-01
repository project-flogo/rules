package rete

type joinTable interface {
	addRow(row joinTableRow) //list of StreamTuples
	getID() int
	len() int
	getMap() map[joinTableRow]joinTableRow
	removeRow(row joinTableRow)
}

type joinTableImpl struct {
	id    int
	table map[joinTableRow]joinTableRow
	idr   []identifier
}

func newJoinTable(identifiers []identifier) joinTable {
	jT := joinTableImpl{}
	jT.initJoinTableImpl(identifiers)
	return &jT
}

func (jt *joinTableImpl) initJoinTableImpl(identifiers []identifier) {
	currentNodeID++
	jt.id = currentNodeID
	jt.idr = identifiers
	jt.table = map[joinTableRow]joinTableRow{}
}

func (jt *joinTableImpl) getID() int {
	return jt.id
}

func (jt *joinTableImpl) addRow(row joinTableRow) {
	jt.table[row] = row
	for i := 0; i < len(row.getHandles()); i++ {
		handle := row.getHandles()[i]
		handle.addJoinTableRowRef(row, jt)
	}
}

func (jt *joinTableImpl) removeRow(row joinTableRow) {
	delete(jt.table, row)
}

func (jt *joinTableImpl) len() int {
	return len(jt.table)
}

func (jt *joinTableImpl) getMap() map[joinTableRow]joinTableRow {
	return jt.table
}
