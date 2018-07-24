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
	OnValueChange(tuple model.StreamTuple)
}

//store any context, may not know all keys upfront
type reteCtxImpl struct {
	cr      conflictRes
	opsList *list.List
	network Network
	rs      model.RuleSession
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

func (rctx *reteCtxImpl) OnValueChange(tuple model.StreamTuple) {
	rctx.opsList.PushBack(newModifyEntry(tuple))
}

func newReteCtxImpl(network Network, rs model.RuleSession) reteCtx {
	reteCtxVal := reteCtxImpl{}
	reteCtxVal.cr = newConflictRes()
	reteCtxVal.opsList = list.New()
	reteCtxVal.network = network
	reteCtxVal.rs = rs
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
