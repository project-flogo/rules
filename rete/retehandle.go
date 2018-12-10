package rete

import (
	"context"

	"github.com/project-flogo/rules/common/model"
)

//Holds a tuple reference and related state
type reteHandle interface {
	nwElemId
	setTuple(tuple model.Tuple)
	getTuple() model.Tuple
	addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable)
	removeJoinTableRowRefs(changedProps map[string]bool)
	removeJoinTable(joinTableID int)
	getTupleKey() model.TupleKey
}

type reteHandleImpl struct {
	nwElemIdImpl
	tuple    model.Tuple
	tupleKey model.TupleKey
	jtRefs   joinTableRefsInHdl
}

func newReteHandleImpl(nw *reteNetworkImpl, tuple model.Tuple) reteHandle {
	h1 := reteHandleImpl{}
	h1.initHandleImpl(nw, tuple)
	return &h1
}

func (hdl *reteHandleImpl) setTuple(tuple model.Tuple) {
	hdl.tuple = tuple
}

func (hdl *reteHandleImpl) initHandleImpl(nw *reteNetworkImpl, tuple model.Tuple) {
	hdl.setID(nw)
	hdl.setTuple(tuple)
	hdl.jtRefs = newJoinTableRefsInHdlImpl()
}

func (hdl *reteHandleImpl) getTuple() model.Tuple {
	return hdl.tuple
}

func (hdl *reteHandleImpl) getTupleKey() model.TupleKey {
	return hdl.tupleKey
}

func getOrCreateHandle(ctx context.Context, tuple model.Tuple) reteHandle {
	reteCtxVar := getReteCtx(ctx)
	return reteCtxVar.getNetwork().getOrCreateHandle(ctx, tuple)
}

func (hdl *reteHandleImpl) addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable) {
	hdl.jtRefs.addEntry(joinTableVar.getID(), joinTableRowVar.getID())
}

func (hdl *reteHandleImpl) removeJoinTableRowRefs(changedProps map[string]bool) {

	tuple := hdl.tuple
	alias := tuple.GetTupleType()

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
		//this can happen if some other handle removed a row as a result of retraction
		if rowIDs == nil {
			continue
		}
		////Remove rows from corresponding join tables
		for e := rowIDs.Front(); e != nil; e = e.Next() {
			rowID := e.Value.(int)
			row := joinTable.removeRow(rowID)
			if row != nil {
				//Remove other refs recursively.
				for _, otherHdl := range row.getHandles() {
					//if otherHdl != nil {
					otherHdl.removeJoinTableRowRefs(nil)
					//}
				}
			}
		}

		//Remove the reference to the table itself
		hdl.jtRefs.removeEntry(joinTableID)
	}
}

//Used when a rule is deleted. See Network.RemoveRule
func (hdl *reteHandleImpl) removeJoinTable(joinTableID int) {
	hdl.jtRefs.removeEntry(joinTableID)
}
