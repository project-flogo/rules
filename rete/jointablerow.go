package rete

type joinTableRow interface {
	nwElemId
	getHandles() []reteHandle
}
