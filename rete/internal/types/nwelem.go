package types

type NwElemId interface {
	SetID(nw Network)
	GetID() int
	GetNw() Network
}

type NwElemIdImpl struct {
	ID int
	Nw Network
}

func (ide *NwElemIdImpl) SetID(nw Network) {
	ide.Nw = nw
	ide.ID = nw.GetIdGenService().GetNextID()
}
func (ide *NwElemIdImpl) GetID() int {
	return ide.ID
}
func (ide *NwElemIdImpl) GetNw() Network {
	return ide.Nw
}
