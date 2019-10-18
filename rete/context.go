package rete

import (
	"container/list"
	"context"
	"fmt"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
)

var reteCTXKEY = model.RetecontextKeyType{}

//store any context, may not know all keys upfront
type reteCtxImpl struct {
	cr      types.ConflictRes
	opsList *list.List
	network *reteNetworkImpl
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

func newReteCtxImpl(network *reteNetworkImpl, rs model.RuleSession) types.ReteCtx {
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

func (rctx *reteCtxImpl) GetConflictResolver() types.ConflictRes {
	return rctx.cr
}

func (rctx *reteCtxImpl) GetOpsList() *list.List {
	return rctx.opsList
}

func (rctx *reteCtxImpl) GetNetwork() types.Network {
	return rctx.network
}

func (rctx *reteCtxImpl) GetRuleSession() model.RuleSession {
	return rctx.rs
}

func (rctx *reteCtxImpl) OnValueChange(tuple model.Tuple, prop string) {

	//if handle does not exist means its new
	if nil != rctx.network.getHandle(context.WithValue(context.Background(), reteCTXKEY, rctx), tuple) {
		rtcModified := rctx.modifyMap[tuple.GetKey().String()]
		if rtcModified == nil {
			rtcModified = NewRtcModified(tuple)
			(rtcModified.(*rtcModifiedImpl)).addProp(prop)
			rctx.modifyMap[tuple.GetKey().String()] = rtcModified
		}
	}
}

func (rctx *reteCtxImpl) GetRtcAdded() map[string]model.Tuple {
	return rctx.addMap
}
func (rctx *reteCtxImpl) GetRtcModified() map[string]model.RtcModified {
	return rctx.rtcModifyMap
}
func (rctx *reteCtxImpl) GetRtcDeleted() map[string]model.Tuple {
	return rctx.deleteMap
}

func (rctx *reteCtxImpl) AddToRtcAdded(tuple model.Tuple) {
	rctx.addMap[tuple.GetKey().String()] = tuple
}

func (rctx *reteCtxImpl) AddToRtcModified(tuple model.Tuple) {
	rctx.addMap[tuple.GetKey().String()] = tuple
}

func (rctx *reteCtxImpl) AddToRtcDeleted(tuple model.Tuple) {
	rctx.deleteMap[tuple.GetKey().String()] = tuple
}

func (rctx *reteCtxImpl) AddRuleModifiedToOpsList() {
	for _, rtcModified := range rctx.modifyMap {
		rctx.GetOpsList().PushBack(newModifyEntry(rtcModified.GetTuple(), rtcModified.GetModifiedProps()))
	}
}

func (rctx *reteCtxImpl) Normalize() {

	//remove from modify map, those in add map
	for k, _ := range rctx.addMap {
		delete(rctx.rtcModifyMap, k)
	}
	//remove from modify map, those in delete map
	for k, _ := range rctx.deleteMap {
		delete(rctx.rtcModifyMap, k)
	}
}

func (rctx *reteCtxImpl) CopyRuleModifiedToRtcModified() {
	for k, v := range rctx.modifyMap {
		rctx.rtcModifyMap[k] = v
	}
}
func (rctx *reteCtxImpl) ResetModified() {
	rctx.modifyMap = make(map[string]model.RtcModified)
}

func (rctx *reteCtxImpl) PrintRtcChangeList() {
	for k := range rctx.GetRtcAdded() {
		fmt.Printf("Added Tuple: [%s]\n", k)
	}
	for k := range rctx.GetRtcModified() {
		fmt.Printf("Modified Tuple: [%s]\n", k)
	}
	for k := range rctx.GetRtcDeleted() {
		fmt.Printf("Deleted Tuple: [%s]\n", k)
	}
}

func getReteCtx(ctx context.Context) types.ReteCtx {
	intr := ctx.Value(reteCTXKEY)
	if intr == nil {
		return nil
	}
	return intr.(types.ReteCtx)
}

func newReteCtx(ctx context.Context, network *reteNetworkImpl, rs model.RuleSession) (context.Context, types.ReteCtx) {
	reteCtxVar := newReteCtxImpl(network, rs)
	ctx = context.WithValue(ctx, reteCTXKEY, reteCtxVar)
	return ctx, reteCtxVar
}

func getOrSetReteCtx(ctx context.Context, network *reteNetworkImpl, rs model.RuleSession) (types.ReteCtx, bool, context.Context) {
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
