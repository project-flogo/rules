package mem

import "github.com/project-flogo/rules/rete/internal/types"

type joinTableRowImpl struct {
	types.NwElemIdImpl
	handles []types.ReteHandle
}

func newJoinTableRow(handles []types.ReteHandle, nw types.Network) types.JoinTableRow {
	jtr := joinTableRowImpl{
		handles: append([]types.ReteHandle{}, handles...),
	}
	jtr.SetID(nw)
	return &jtr
}

func (jtr *joinTableRowImpl) Write() {

}

func (jtr *joinTableRowImpl) GetHandles() []types.ReteHandle {
	return jtr.handles
}
