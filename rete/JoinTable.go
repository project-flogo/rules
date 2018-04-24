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
	idr   []Identifier
}

func newJoinTable(identifiers []Identifier) joinTable {
	jT := joinTableImpl{}
	jT.initJoinTableImpl(identifiers)
	return &jT
}

func (joinTableImplVar *joinTableImpl) initJoinTableImpl(identifiers []Identifier) {
	currentNodeID++
	joinTableImplVar.id = currentNodeID
	joinTableImplVar.idr = identifiers
	joinTableImplVar.table = map[joinTableRow]joinTableRow{}
}

func (joinTableImplVar *joinTableImpl) getID() int {
	return joinTableImplVar.id
}

func (joinTableImplVar *joinTableImpl) addRow(row joinTableRow) {
	joinTableImplVar.table[row] = row
	for i := 0; i < len(row.getHandles()); i++ {
		handle := row.getHandles()[i]
		handle.addJoinTableRowRef(row, joinTableImplVar)
	}
}

func (joinTableImplVar *joinTableImpl) removeRow(row joinTableRow) {
	delete(joinTableImplVar.table, row)
}

func (joinTableImplVar *joinTableImpl) len() int {
	return len(joinTableImplVar.table)
}

func (joinTableImplVar *joinTableImpl) getMap() map[joinTableRow]joinTableRow {
	return joinTableImplVar.table
}
