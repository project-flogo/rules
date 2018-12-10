package rete

type joinTableRow interface {
	nwElemId
	getHandles() []reteHandle
}

type joinTableRowImpl struct {
	nwElemIdImpl
	handles []reteHandle
}

func newJoinTableRow(handles []reteHandle, nw Network) joinTableRow {
	jtr := joinTableRowImpl{}
	jtr.initJoinTableRow(handles, nw)
	return &jtr
}

func (jtr *joinTableRowImpl) initJoinTableRow(handles []reteHandle, nw Network) {
	jtr.setID(nw)
	jtr.handles = append([]reteHandle{}, handles...)
}

func (jtr *joinTableRowImpl) getHandles() []reteHandle {
	return jtr.handles
}
