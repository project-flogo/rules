package rete

import (
	"context"

	"github.com/project-flogo/rules/common/model"
)

//Holds a tuple reference and related state
type reteHandle interface {
	setTuple(tuple model.Tuple)
	getTuple() model.Tuple
	addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable)
	removeJoinTableRowRefs(changedProps map[string]bool)
	removeJoinTable(joinTableID int)
}

type reteHandleImpl struct {
	tuple     model.Tuple
	rtcStatus uint8
	nw        Network
	rhRef     reteHandleRefs
}

func (hdl *reteHandleImpl) setTuple(tuple model.Tuple) {
	hdl.tuple = tuple
}

func (hdl *reteHandleImpl) initHandleImpl() {
	hdl.rhRef = newReteHandleRefsImpl()
}

func (hdl *reteHandleImpl) getTuple() model.Tuple {
	return hdl.tuple
}

func getOrCreateHandle(ctx context.Context, tuple model.Tuple) reteHandle {
	reteCtxVar := getReteCtx(ctx)
	return reteCtxVar.getNetwork().getOrCreateHandle(ctx, tuple)
}

func (hdl *reteHandleImpl) addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable) {
	hdl.rhRef.addEntry(joinTableVar.getID(), joinTableRowVar.getID())
}

func (hdl *reteHandleImpl) removeJoinTableRowRefs(changedProps map[string]bool) {

	tuple := hdl.tuple
	alias := tuple.GetTupleType()

	//emptyJoinTables := list.New()

	hdlTblIter := hdl.newHdlTblIterator()

	for hdlTblIter.hasNext() {
		joinTableID, rowIDs := hdlTblIter.next()
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

		hdl.rhRef.removeEntry(joinTableID)
	}
}

//Used when a rule is deleted. See Network.RemoveRule
func (hdl *reteHandleImpl) removeJoinTable(joinTableID int) {
	hdl.rhRef.removeEntry(joinTableID)
}

//func (hdl *reteHandleImpl) deleteRefsToJoinTables (jointTableID int) {
//	delete (hdl.tablesAndRows, jointTableID)
//}
