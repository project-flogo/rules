package rete

import (
	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/utils"
)

//Holds a stream tuple reference and related state
type reteHandle interface {
	setTuple(streamTuple model.StreamTuple)
	getTuple() model.StreamTuple
	addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable)
	removeJoinTableRowRefs()
	removeJoinTable(joinTableVar joinTable)
}

type handleImpl struct {
	tuple         model.StreamTuple
	tablesAndRows map[joinTable]utils.ArrayList
}

func (handleImplVar *handleImpl) setTuple(tuple model.StreamTuple) {
	handleImplVar.tuple = tuple
}

func (handleImplVar *handleImpl) initHandleImpl() {
	handleImplVar.tablesAndRows = make(map[joinTable]utils.ArrayList)
}

func (handleImplVar *handleImpl) getTuple() model.StreamTuple {
	return handleImplVar.tuple
}

func getOrCreateHandle(tuple model.StreamTuple) reteHandle {
	h := allHandles[tuple]
	if h == nil {
		h1 := handleImpl{}
		h1.initHandleImpl()
		h1.setTuple(tuple)
		h = &h1
		allHandles[tuple] = h
	}
	return h
}

func (handleImplVar *handleImpl) addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable) {

	rowsForJoinTable := handleImplVar.tablesAndRows[joinTableVar]
	if rowsForJoinTable == nil {
		rowsForJoinTable = utils.NewArrayList()
		handleImplVar.tablesAndRows[joinTableVar] = rowsForJoinTable
	}
	rowsForJoinTable.Add(joinTableRowVar)

}

func (handleImplVar *handleImpl) removeJoinTableRowRefs() {

	emptyJoinTables := utils.NewArrayList()

	for joinTable, listOfRows := range handleImplVar.tablesAndRows {
		for i := 0; i < listOfRows.Len(); i++ {
			row := listOfRows.Get(i).(joinTableRow)
			joinTable.removeRow(row)
		}
		if joinTable.len() == 0 {
			emptyJoinTables.Add(joinTable)
		}
	}

	for i := 0; i < emptyJoinTables.Len(); i++ {
		emptyJoinTable := emptyJoinTables.Get(i).(joinTable)
		delete(handleImplVar.tablesAndRows, emptyJoinTable)
	}
}

//Used when a rule is deleted. See Network.RemoveRule
func (handleImplVar *handleImpl) removeJoinTable(joinTableVar joinTable) {
	_, ok := handleImplVar.tablesAndRows[joinTableVar]
	if ok {
		delete(handleImplVar.tablesAndRows, joinTableVar)
	}
}
