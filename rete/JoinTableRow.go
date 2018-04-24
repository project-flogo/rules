package rete

type joinTableRow interface {
	getHandles() []reteHandle
}

type joinTableRowImpl struct {
	handles []reteHandle
}

func newJoinTableRow(handles []reteHandle) joinTableRow {
	joinTableRowImplVar := joinTableRowImpl{}
	joinTableRowImplVar.initJoinTableRow(handles)
	return &joinTableRowImplVar
}

func (joinTableRowImplVar *joinTableRowImpl) initJoinTableRow(handles []reteHandle) {
	joinTableRowImplVar.handles = append([]reteHandle{}, handles...)
}

func (joinTableRowImplVar *joinTableRowImpl) getHandles() []reteHandle {
	return joinTableRowImplVar.handles
}
