package rete

type joinTableRow interface {
	getID() int
	getHandles() []reteHandle
}

type joinTableRowImpl struct {
	id      int
	handles []reteHandle
}

func newJoinTableRow(handles []reteHandle, id int) joinTableRow {
	jtr := joinTableRowImpl{}
	jtr.initJoinTableRow(handles, id)
	return &jtr
}

func (jtr *joinTableRowImpl) initJoinTableRow(handles []reteHandle, id int) {
	jtr.handles = append([]reteHandle{}, handles...)
	jtr.id = id
}

func (jtr *joinTableRowImpl) getHandles() []reteHandle {
	return jtr.handles
}

func (jtr *joinTableRowImpl) getID() int {
	return jtr.id
}
