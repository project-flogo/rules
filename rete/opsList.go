package rete

import (
	"context"

	"github.com/TIBCOSoftware/bego/common/model"
)

type opsEntry interface {
	execute(ctx context.Context)
}

type opsEntryImpl struct {
	tuple model.StreamTuple
}

//Assert entry

type assertEntry interface {
	opsEntry
}

type assertEntryImpl struct {
	opsEntryImpl
}

func newAssertEntry(tuple model.StreamTuple) assertEntry {
	aEntry := assertEntryImpl{}
	aEntry.tuple = tuple
	return &aEntry
}

func (ai *assertEntryImpl) execute(ctx context.Context) {
	reteCtx := getReteCtx(ctx)
	reteCtx.getNetwork().assertInternal(ctx, ai.tuple)
}

//Modify Entry

type modifyEntry interface {
	opsEntry
}

type modifyEntryImpl struct {
	opsEntryImpl
}

func newModifyEntry(tuple model.StreamTuple) modifyEntry {
	mEntry := modifyEntryImpl{}
	mEntry.tuple = tuple
	return &mEntry
}

func (me *modifyEntryImpl) execute(ctx context.Context) {
	reteCtx := getReteCtx(ctx)
	reteCtx.getConflictResolver().deleteAgendaFor(ctx, me.tuple)
	reteCtx.getNetwork().Retract(me.tuple)
	reteCtx.getNetwork().Assert(ctx, me.tuple)
}

//Delete Entry

type deleteEntry interface {
	opsEntry
}

type deleteEntryImpl struct {
	opsEntryImpl
}

func newDeleteEntry(tuple model.StreamTuple) deleteEntry {
	dEntry := deleteEntryImpl{}
	dEntry.tuple = tuple
	return &dEntry
}

func (de *deleteEntryImpl) execute(ctx context.Context) {
	reteCtx := getReteCtx(ctx)
	reteCtx.getConflictResolver().deleteAgendaFor(ctx, de.tuple)
	reteCtx.getNetwork().Retract(de.tuple)
}
