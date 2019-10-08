package rete

import (
	"container/list"
	"context"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
)

type conflictResImpl struct {
	agendaList list.List
}

func newConflictRes() types.ConflictRes {
	cr := conflictResImpl{}
	cr.initCR()
	return &cr
}

func (cr *conflictResImpl) initCR() {
	cr.agendaList = list.List{}
}

func (cr *conflictResImpl) AddAgendaItem(rule model.Rule, tupleMap map[model.TupleType]model.Tuple) {
	item := newAgendaItem(rule, tupleMap)
	v := rule.GetPriority()
	found := false
	for e := cr.agendaList.Front(); e != nil; e = e.Next() {
		curr := e.Value.(agendaItem)
		if v <= curr.getRule().GetPriority() {
			cr.agendaList.InsertBefore(item, e)
			found = true
			break
		}
	}
	if !found {
		cr.agendaList.PushBack(item)
	}
}

func (cr *conflictResImpl) ResolveConflict(ctx context.Context) {
	var item agendaItem

	front := cr.agendaList.Front()
	for front != nil {
		val := cr.agendaList.Remove(front)
		if val != nil {
			item = val.(agendaItem)
			actionTuples := item.getTuples()
			// execute rule action service
			aService := item.getRule().GetActionService()
			if aService != nil {
				reteCtxV := getReteCtx(ctx)
				aService.Execute(ctx, reteCtxV.GetRuleSession(), item.getRule().GetName(), actionTuples, item.getRule().GetContext())
			}
		}

		reteCtxV := getReteCtx(ctx)

		reteCtxV.AddRuleModifiedToOpsList()

		reteCtxV.CopyRuleModifiedToRtcModified()
		//action scoped, clear it for the next action
		reteCtxV.ResetModified()

		if reteCtxV != nil {
			opsFront := reteCtxV.GetOpsList().Front()
			for opsFront != nil {
				opsVal := reteCtxV.GetOpsList().Remove(opsFront)
				oprn := opsVal.(opsEntry)
				oprn.execute(ctx)
				opsFront = reteCtxV.GetOpsList().Front()
			}
		}

		front = cr.agendaList.Front()
	}

	reteCtxV := getReteCtx(ctx)
	reteCtxV.Normalize()
	//reteCtxV.printRtcChangeList()

}

func (cr *conflictResImpl) DeleteAgendaFor(ctx context.Context, modifiedTuple model.Tuple, changeProps map[string]bool) {

	hdlModified, _ := getOrCreateHandle(ctx, modifiedTuple)

	for e := cr.agendaList.Front(); e != nil; {
		item := e.Value.(agendaItem)
		next := e.Next()
		for _, tuple := range item.getTuples() {
			hdl, _ := getOrCreateHandle(ctx, tuple)
			if hdl == hdlModified { //this agendaitem has the modified tuple, remove the agenda item!
				toRemove := true
				//check if the rule depends on this change prop
				if changeProps != nil {
					if depProps, found := item.getRule().GetDeps()[tuple.GetTupleType()]; found {
						if len(depProps) > 0 {
							for prop, _ := range depProps {
								if _, fnd := changeProps[prop]; fnd {
									toRemove = true
									break
								}
							}
						} else {
							toRemove = false
						}
					} else {
						toRemove = false
					}
				}
				if toRemove {
					cr.agendaList.Remove(e)
					break
				}
			}
		}
		e = next
	}

}
