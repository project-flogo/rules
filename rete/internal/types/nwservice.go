package types

import "github.com/project-flogo/rules/common/model"

type NwService interface {
	model.Service
	GetNw() Network
}

type NwServiceImpl struct {
	Nw Network
}

func (nws *NwServiceImpl) GetNw() Network {
	return nws.Nw
}
