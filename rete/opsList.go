package rete

import (
	"context"

	"github.com/project-flogo/rules/common/model"
)

type opsEntry interface {
	execute(ctx context.Context)
}

type opsEntryImpl struct {
	tuple       model.Tuple
	changeProps map[string]bool
}

//Assert entry

type assertEntry interface {
	opsEntry
}

type assertEntryImpl struct {
	opsEntryImpl
	mode RtcOprn
}

func newAssertEntry(tuple model.Tuple, changeProps map[string]bool, mode RtcOprn) assertEntry {
	aEntry := assertEntryImpl{}
	aEntry.tuple = tuple
	aEntry.changeProps = changeProps
	aEntry.mode = mode
	return &aEntry
}

func (ai *assertEntryImpl) execute(ctx context.Context) {
	reteCtx := getReteCtx(ctx)
	reteCtx.getNetwork().assertInternal(ctx, ai.tuple, ai.changeProps, ai.mode, "")
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
	reteCtx.getConflictResolver().deleteAgendaFor(ctx, me.tuple, me.changeProps)
	reteCtx.getNetwork().Retract(ctx, me.tuple, me.changeProps, MODIFY)
	reteCtx.getNetwork().Assert(ctx, reteCtx.getRuleSession(), me.tuple, me.changeProps, MODIFY)
}

//Delete Entry

type deleteEntry interface {
	opsEntry
}

type deleteEntryImpl struct {
	opsEntryImpl
	mode RtcOprn
}

func newDeleteEntry(tuple model.Tuple, mode RtcOprn, changeProps map[string]bool) deleteEntry {
	dEntry := deleteEntryImpl{}
	dEntry.tuple = tuple
	dEntry.mode = mode
	dEntry.changeProps = changeProps
	return &dEntry
}

func (de *deleteEntryImpl) execute(ctx context.Context) {
	reteCtx := getReteCtx(ctx)
	reteCtx.getConflictResolver().deleteAgendaFor(ctx, de.tuple, de.changeProps)
	reteCtx.getNetwork().retractInternal(ctx, de.tuple, de.changeProps, de.mode)
}
