package rete

import (
	"container/list"
	"context"
)

type retecontextKeyType struct {
}

var reteCTXKEY = retecontextKeyType{}

type reteCtx interface {
	getConflictResolver() conflictRes
	getOpsList() *list.List
	getNetwork() Network
}

//store any context, may not know all keys upfront
type reteCtxImpl struct {
	cr      conflictRes
	opsList *list.List
	network Network
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

func newReteCtxImpl(network Network) reteCtx {
	reteCtxVal := reteCtxImpl{}
	reteCtxVal.cr = newConflictRes()
	reteCtxVal.opsList = list.New()
	reteCtxVal.network = network
	return &reteCtxVal
}

func getReteCtx(ctx context.Context) reteCtx {
	intr := ctx.Value(reteCTXKEY)
	if intr == nil {
		return nil
	} else {
		return intr.(reteCtx)
	}
}

func newCtx(network Network) (context.Context, reteCtx) {
	reteCtxVar := newReteCtxImpl(network)
	ctx := context.WithValue(context.Background(), reteCTXKEY, reteCtxVar)
	return ctx, reteCtxVar
}

func newReteCtx(ctx context.Context, network Network) (context.Context, reteCtx) {
	reteCtxVar := newReteCtxImpl(network)
	ctx = context.WithValue(ctx, reteCTXKEY, reteCtxVar)
	return ctx, reteCtxVar
}

func getOrSetReteCtx(ctx context.Context, network Network) (reteCtx, bool) {
	found := false
	reteCtxVar := getReteCtx(ctx)
	if reteCtxVar == nil {
		_, reteCtxVar = newReteCtx(ctx, network)
		found = false
	} else {
		found = true
	}
	return reteCtxVar, found
}
