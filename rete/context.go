package rete

import (
	"container/list"
	"context"

	"github.com/TIBCOSoftware/bego/common/model"
)

var reteCTXKEY = model.RetecontextKeyType{}

type reteCtx interface {
	getConflictResolver() conflictRes
	getOpsList() *list.List
	getNetwork() Network
	getRuleSession() model.RuleSession
	OnValueChange(tuple model.StreamTuple, prop string)
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
	addMap map[model.StreamTuple]bool

	//in each action, this map is updated with modifications
	//key is the added tuple, value is a map of the changed property to true
	// (we simply want unique modified tuples and unique props in each that changed)
	modifyMap map[model.StreamTuple]map[string]bool

	//in each action, this map is updated with deletions
	//key is the deleted tuple value is always true
	// (we simply want a unique set of deleted tuples)
	deleteMap map[model.StreamTuple]bool
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

func (rctx *reteCtxImpl) OnValueChange(tuple model.StreamTuple, prop string) {
	propMap := rctx.modifyMap[tuple]
	if propMap == nil {
		propMap = make(map[string]bool)
		propMap[prop] = true
		rctx.modifyMap[tuple] = propMap
	} else {
		propMap[prop] = true
	}
}

func newReteCtxImpl(network Network, rs model.RuleSession) reteCtx {
	reteCtxVal := reteCtxImpl{}
	reteCtxVal.cr = newConflictRes()
	reteCtxVal.opsList = list.New()
	reteCtxVal.network = network
	reteCtxVal.rs = rs
	reteCtxVal.addMap = make (map[model.StreamTuple]bool)
	reteCtxVal.modifyMap = make(map[model.StreamTuple]map[string]bool)
	reteCtxVal.deleteMap = make (map[model.StreamTuple]bool)
	return &reteCtxVal
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
