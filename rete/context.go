package rete

import (
	"container/list"
	"context"

	"fmt"
	"github.com/project-flogo/rules/common/model"
)

var reteCTXKEY = model.RetecontextKeyType{}

type reteCtx interface {
	getConflictResolver() conflictRes
	getOpsList() *list.List
	getNetwork() Network
	getRuleSession() model.RuleSession
	OnValueChange(tuple model.Tuple, prop string)

	getRtcAdded() map[string]model.Tuple
	getRtcModified() map[string]model.RtcModified
	getRtcDeleted() map[string]model.Tuple

	addToRtcAdded(tuple model.Tuple)
	addToRtcModified(tuple model.Tuple)
	addToRtcDeleted(tuple model.Tuple)
	addRuleModifiedToOpsList()

	normalize()
	copyRuleModifiedToRtcModified()
	resetModified()

	printRtcChangeList()
}

//store any context, may not know all keys upfront
type reteCtxImpl struct {
	cr      conflictRes
	opsList *list.List
	network Network
	rs      model.RuleSession

	//newly added tuples in the current RTC
	addMap map[string]model.Tuple

	//modified tuples in the current rule action
	modifyMap map[string]model.RtcModified

	//deleted (which is different than simply retracted) tuples in the current RTC (
	deleteMap map[string]model.Tuple

	//modified tuples in the current RTC
	rtcModifyMap map[string]model.RtcModified
}

func newReteCtxImpl(network Network, rs model.RuleSession) reteCtx {
	reteCtxVal := reteCtxImpl{}
	reteCtxVal.cr = newConflictRes()
	reteCtxVal.opsList = list.New()
	reteCtxVal.network = network
	reteCtxVal.rs = rs
	reteCtxVal.addMap = make(map[string]model.Tuple)
	reteCtxVal.modifyMap = make(map[string]model.RtcModified)
	reteCtxVal.rtcModifyMap = make(map[string]model.RtcModified)
	reteCtxVal.deleteMap = make(map[string]model.Tuple)
	return &reteCtxVal
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
		rtcModified := rctx.modifyMap[tuple.GetKey().String()]
		if rtcModified == nil {
			rtcModified = NewRtcModified(tuple)
			(rtcModified.(*rtcModifiedImpl)).addProp(prop)
			rctx.modifyMap[tuple.GetKey().String()] = rtcModified
		}
	}
}

func (rctx *reteCtxImpl) getRtcAdded() map[string]model.Tuple {
	return rctx.addMap
}
func (rctx *reteCtxImpl) getRtcModified() map[string]model.RtcModified {
	return rctx.rtcModifyMap
}
func (rctx *reteCtxImpl) getRtcDeleted() map[string]model.Tuple {
	return rctx.deleteMap
}

func (rctx *reteCtxImpl) addToRtcAdded(tuple model.Tuple) {
	rctx.addMap[tuple.GetKey().String()] = tuple
}

func (rctx *reteCtxImpl) addToRtcModified(tuple model.Tuple) {
	rctx.addMap[tuple.GetKey().String()] = tuple
}

func (rctx *reteCtxImpl) addToRtcDeleted(tuple model.Tuple) {
	rctx.deleteMap[tuple.GetKey().String()] = tuple
}

func (rctx *reteCtxImpl) addRuleModifiedToOpsList() {
	for _, rtcModified := range rctx.modifyMap {
		rctx.getOpsList().PushBack(newModifyEntry(rtcModified.GetTuple(), rtcModified.GetModifiedProps()))
	}
}

func (rctx *reteCtxImpl) normalize() {

	//remove from modify map, those in add map
	for k, _ := range rctx.addMap {
		delete(rctx.rtcModifyMap, k)
	}
	//remove from modify map, those in delete map
	for k, _ := range rctx.deleteMap {
		delete(rctx.rtcModifyMap, k)
	}
}

func (rctx *reteCtxImpl) copyRuleModifiedToRtcModified() {
	for k, v := range rctx.modifyMap {
		rctx.rtcModifyMap[k] = v
	}
}
func (rctx *reteCtxImpl) resetModified() {
	rctx.modifyMap = make(map[string]model.RtcModified)
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
