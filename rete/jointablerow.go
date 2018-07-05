package rete

type joinTableRow interface {
	getHandles() []reteHandle
}

type joinTableRowImpl struct {
	handles []reteHandle
}

func newJoinTableRow(handles []reteHandle) joinTableRow {
	jtr := joinTableRowImpl{}
	jtr.initJoinTableRow(handles)
	return &jtr
}

func (jtr *joinTableRowImpl) initJoinTableRow(handles []reteHandle) {
	jtr.handles = append([]reteHandle{}, handles...)
}

func (jtr *joinTableRowImpl) getHandles() []reteHandle {
	return jtr.handles
}
