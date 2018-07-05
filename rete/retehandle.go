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

func (hdl *handleImpl) setTuple(tuple model.StreamTuple) {
	hdl.tuple = tuple
}

func (hdl *handleImpl) initHandleImpl() {
	hdl.tablesAndRows = make(map[joinTable]utils.ArrayList)
}

func (hdl *handleImpl) getTuple() model.StreamTuple {
	return hdl.tuple
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

func (hdl *handleImpl) addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable) {

	rowsForJoinTable := hdl.tablesAndRows[joinTableVar]
	if rowsForJoinTable == nil {
		rowsForJoinTable = utils.NewArrayList()
		hdl.tablesAndRows[joinTableVar] = rowsForJoinTable
	}
	rowsForJoinTable.Add(joinTableRowVar)

}

func (hdl *handleImpl) removeJoinTableRowRefs() {

	emptyJoinTables := utils.NewArrayList()

	for joinTable, listOfRows := range hdl.tablesAndRows {
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
		delete(hdl.tablesAndRows, emptyJoinTable)
	}
}

//Used when a rule is deleted. See Network.RemoveRule
func (hdl *handleImpl) removeJoinTable(joinTableVar joinTable) {
	_, ok := hdl.tablesAndRows[joinTableVar]
	if ok {
		delete(hdl.tablesAndRows, joinTableVar)
	}
}
