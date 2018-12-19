package mem

import "github.com/project-flogo/rules/rete/internal/types"

type joinTableRowImpl struct {
	types.NwElemIdImpl
	handles []types.ReteHandle
}

func newJoinTableRow(handles []types.ReteHandle, nw types.Network) types.JoinTableRow {
	jtr := joinTableRowImpl{}
	jtr.initJoinTableRow(handles, nw)
	return &jtr
}

func (jtr *joinTableRowImpl) initJoinTableRow(handles []types.ReteHandle, nw types.Network) {
	jtr.SetID(nw)
	jtr.handles = append([]types.ReteHandle{}, handles...)
}

func (jtr *joinTableRowImpl) GetHandles() []types.ReteHandle {
	return jtr.handles
}
