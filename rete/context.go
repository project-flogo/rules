package rete

import (
	"container/list"
	"context"

	"github.com/project-flogo/rules/common/model"
	"fmt"
)

var reteCTXKEY = model.RetecontextKeyType{}

type reteCtx interface {
	getConflictResolver() conflictRes
	getOpsList() *list.List
	getNetwork() Network
	getRuleSession() model.RuleSession
	OnValueChange(tuple model.Tuple, prop string)

	getRtcAdded() map[string]model.Tuple
	getRtcModified() map[string]map[string]bool
	getRtcDeleted() map[string]model.Tuple

	addToRtcAdded(tuple model.Tuple)
	resetModified()
	addToRtcDeleted(tuple model.Tuple)
	normalize()

	addToRtcModified(tuple model.Tuple)
	addRuleModifiedToOpsList()
	copyRuleModifiedToRtcModified()

	printRtcChangeList()

}

//store any context, may not know all keys upfront
type reteCtxImpl struct {
	cr      conflictRes
	opsList *list.List
	network Network
	rs      model.RuleSession

	//in each action, this map is updated with new ones
	//key is the added tuple, value is true
	// (we simply want the unique set of newly added tuples)
	addMap map[string]model.Tuple

	//in each action, this map is updated with modifications
	//key is the added tuple, value is a map of the changed property to true
	// (we simply want unique modified tuples and unique props in each that changed)
	modifyMap map[string]map[string]bool

	//in each action, this map is updated with deletions
	//key is the deleted tuple value is always true
	// (we simply want a unique set of deleted tuples)
	deleteMap map[string]model.Tuple

	//in each action, this map is updated with modifications
	//key is the added tuple, value is a map of the changed property to true
	// (we simply want unique modified tuples and unique props in each that changed)
	rtcModifyMap map[string]map[string]bool

}

func (rctx *reteCtxImpl) getConflictResolver() conflictRes {
	return rctx.cr
}

func (rctx *reteCtxImpl) getOpsList() *list.List {
	return rctx.opsList
}

func (rctx *reteCtxImpl) getNetwork() Network {
	return rctx.network
}

func (rctx *reteCtxImpl) getRuleSession() model.RuleSession {
	return rctx.rs
}

func (rctx *reteCtxImpl) OnValueChange(tuple model.Tuple, prop string) {

	//if handle does not exist means its new
	if nil != rctx.network.getHandle(tuple) {
		propMap := rctx.modifyMap[tuple.GetKey().String()]
		if propMap == nil {
			propMap = make(map[string]bool)
			rctx.modifyMap[tuple.GetKey().String()] = propMap
		}
		propMap[prop] = true
	}
}

func (rctx *reteCtxImpl) normalize() {

	//remove from modify map, those in add map
	for k, _ := range rctx.addMap {
		delete(rctx.modifyMap, k)
	}
	//remove from modify map, those in delete map
	for k, _ := range rctx.deleteMap {
		delete(rctx.modifyMap, k)
	}
}

func (rctx *reteCtxImpl) getRtcAdded() map[string]model.Tuple {
	return rctx.addMap
}
func (rctx *reteCtxImpl) getRtcModified() map[string]map[string]bool {
	return rctx.rtcModifyMap
}
func (rctx *reteCtxImpl) getRtcDeleted() map[string]model.Tuple {
	return rctx.deleteMap
}

func newReteCtxImpl(network Network, rs model.RuleSession) reteCtx {
	reteCtxVal := reteCtxImpl{}
	reteCtxVal.cr = newConflictRes()
	reteCtxVal.opsList = list.New()
	reteCtxVal.network = network
	reteCtxVal.rs = rs
	reteCtxVal.addMap = make(map[string]model.Tuple)
	reteCtxVal.modifyMap = make(map[string]map[string]bool)
	reteCtxVal.rtcModifyMap = make(map[string]map[string]bool)
	reteCtxVal.deleteMap = make(map[string]model.Tuple)
	return &reteCtxVal
}

func (rctx *reteCtxImpl) addToRtcAdded (tuple model.Tuple) {
	rctx.addMap[tuple.GetKey().String()] = tuple
}

func (rctx *reteCtxImpl) addToRtcModified (tuple model.Tuple) {
	rctx.addMap[tuple.GetKey().String()] = tuple
}

func (rctx *reteCtxImpl) addRuleModifiedToOpsList() {
	for tupleKey, props := range rctx.modifyMap {
		rctx.getOpsList().PushBack(newModifyEntry(rctx.getNetwork().GetAssertedTupleByStringKey(tupleKey), props))
	}
}
func (rctx *reteCtxImpl) copyRuleModifiedToRtcModified () {
	for k, v := range rctx.modifyMap {
		rctx.rtcModifyMap[k] = v
	}
}
func (rctx *reteCtxImpl) resetModified() {
	rctx.modifyMap = make(map[string]map[string]bool)
}

func (rctx *reteCtxImpl) addToRtcDeleted (tuple model.Tuple) {
	rctx.deleteMap[tuple.GetKey().String()] = tuple
}

func (rctx *reteCtxImpl) printRtcChangeList() {
	for k, _ := range rctx.getRtcAdded() {

		fmt.Printf("Added Tuple: [%s]\n", k)

	}
	for k, _ := range rctx.getRtcModified() {

		fmt.Printf("Modified Tuple: [%s]\n", k)

	}
	for k, _ := range rctx.getRtcDeleted() {

		fmt.Printf("Deleted Tuple: [%s]\n", k)

	}

}

func getReteCtx(ctx context.Context) reteCtx {
	intr := ctx.Value(reteCTXKEY)
	if intr == nil {
		return nil
	}
	return intr.(reteCtx)
}

func newReteCtx(ctx context.Context, network Network, rs model.RuleSession) (context.Context, reteCtx) {
	reteCtxVar := newReteCtxImpl(network, rs)
	ctx = context.WithValue(ctx, reteCTXKEY, reteCtxVar)
	return ctx, reteCtxVar
}

func getOrSetReteCtx(ctx context.Context, network Network, rs model.RuleSession) (reteCtx, bool, context.Context) {
	isRecursive := false
	newCtx := ctx
	reteCtxVar := getReteCtx(ctx)
	if reteCtxVar == nil {
		newCtx, reteCtxVar = newReteCtx(ctx, network, rs)
		isRecursive = false
	} else {
		isRecursive = true
	}
	return reteCtxVar, isRecursive, newCtx
}
