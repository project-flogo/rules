package rete

import (
	"context"

	"github.com/TIBCOSoftware/bego/common/model"
)

type opsEntry interface {
	execute(ctx context.Context)
}

type opsEntryImpl struct {
	tuple model.Tuple
	changeProps map[string]bool
}

//Assert entry

type assertEntry interface {
	opsEntry
}

type assertEntryImpl struct {
	opsEntryImpl

}

func newAssertEntry(tuple model.Tuple, changeProps map[string]bool) assertEntry {
	aEntry := assertEntryImpl{}
	aEntry.tuple = tuple
	aEntry.changeProps = changeProps
	return &aEntry
}

func (ai *assertEntryImpl) execute(ctx context.Context) {
	reteCtx := getReteCtx(ctx)
	reteCtx.getNetwork().assertInternal(ctx, ai.tuple, ai.changeProps)
}

//Modify Entry

type modifyEntry interface {
	opsEntry
}

type modifyEntryImpl struct {
	opsEntryImpl
}

func newModifyEntry(tuple model.Tuple, changeProps map[string]bool) modifyEntry {
	mEntry := modifyEntryImpl{}
	mEntry.tuple = tuple
	mEntry.changeProps = changeProps
	return &mEntry
}

func (me *modifyEntryImpl) execute(ctx context.Context) {
	reteCtx := getReteCtx(ctx)
	reteCtx.getConflictResolver().deleteAgendaFor(ctx, me.tuple)
	reteCtx.getNetwork().Retract(ctx, me.tuple, me.changeProps)
	reteCtx.getNetwork().Assert(ctx, reteCtx.getRuleSession(), me.tuple, me.changeProps)
}

//Delete Entry

type deleteEntry interface {
	opsEntry
}

type deleteEntryImpl struct {
	opsEntryImpl
}

func newDeleteEntry(tuple model.Tuple) deleteEntry {
	dEntry := deleteEntryImpl{}
	dEntry.tuple = tuple
	return &dEntry
}

func (de *deleteEntryImpl) execute(ctx context.Context) {
	reteCtx := getReteCtx(ctx)
	reteCtx.getConflictResolver().deleteAgendaFor(ctx, de.tuple)
	reteCtx.getNetwork().Retract(ctx, de.tuple, de.changeProps)
}
