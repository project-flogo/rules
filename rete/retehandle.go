package rete

import (
	"container/list"
	"context"

	"github.com/project-flogo/rules/common/model"
)

//Holds a tuple reference and related state
type reteHandle interface {
	setTuple(tuple model.Tuple)
	getTuple() model.Tuple
	addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable)
	removeJoinTableRowRefs(changedProps map[string]bool)
	removeJoinTable(joinTableVar joinTable)
}

type handleImpl struct {
	tuple         model.Tuple
	//keys are jointable-ids and values are lists of row-ids in the corresponding join table
	tablesAndRows map[int]*list.List

	rtcStatus uint8
	nw Network
}

func (hdl *handleImpl) setTuple(tuple model.Tuple) {
	hdl.tuple = tuple
}

func (hdl *handleImpl) initHandleImpl() {
	hdl.tablesAndRows = make(map[int]*list.List)
	hdl.rtcStatus = 0x00
}

func (hdl *handleImpl) getTuple() model.Tuple {
	return hdl.tuple
}

func getOrCreateHandle(ctx context.Context, tuple model.Tuple) reteHandle {
	reteCtxVar := getReteCtx(ctx)
	return reteCtxVar.getNetwork().getOrCreateHandle(ctx, tuple)
}

func (hdl *handleImpl) addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable) {

	rowsForJoinTable := hdl.tablesAndRows[joinTableVar.getID()]
	if rowsForJoinTable == nil {
		rowsForJoinTable = list.New()
		hdl.tablesAndRows[joinTableVar.getID()] = rowsForJoinTable
	}
	rowsForJoinTable.PushBack(joinTableRowVar.getID())

}

func (hdl *handleImpl) removeJoinTableRowRefs(changedProps map[string]bool) {

	tuple := hdl.tuple
	alias := tuple.GetTupleType()

	emptyJoinTables := list.New()

	for joinTableID, rowIDs := range hdl.tablesAndRows {
		joinTable := hdl.nw.getJoinTable(joinTableID)
		toDelete := false
		if changedProps != nil {
			rule := joinTable.getRule()
			depProps, found := rule.GetDeps()[alias]
			if found { // rule depends on this type
				for changedProp := range changedProps {
					_, foundProp := depProps[changedProp]
					if foundProp {
						toDelete = true
						break
					}
				}
			}
		} else {
			toDelete = true
		}

		if !toDelete {
			continue
		}

		for e := rowIDs.Front(); e != nil; e = e.Next() {
			rowID := e.Value.(int)
			joinTable.removeRow(rowID)
		}
		if joinTable.len() == 0 {
			emptyJoinTables.PushBack(joinTable.getID())
		}
	}

	for e := emptyJoinTables.Front(); e != nil; e = e.Next() {
		joinTableID := e.Value.(int)
		delete(hdl.tablesAndRows, joinTableID)
	}
}

//Used when a rule is deleted. See Network.RemoveRule
func (hdl *handleImpl) removeJoinTable(joinTableVar joinTable) {
	_, ok := hdl.tablesAndRows[joinTableVar.getID()]
	if ok {
		delete(hdl.tablesAndRows, joinTableVar.getID())
	}
}
