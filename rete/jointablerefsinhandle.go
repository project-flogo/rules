package rete

type joinTableRefsInHdl interface {
	addEntry(jointTableID int, rowID int)
	removeEntry(jointTableID int)
}
