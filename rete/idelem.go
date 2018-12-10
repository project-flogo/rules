package rete

type nwElemId interface {
	setID(nw Network)
	getID() int
}
type nwElemIdImpl struct {
	ID int
	nw Network
}

func (ide *nwElemIdImpl) setID(nw Network) {
	ide.nw = nw
	ide.ID = nw.incrementAndGetId()
}
func (ide *nwElemIdImpl) getID() int {
	return ide.ID
}
