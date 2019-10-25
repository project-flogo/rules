package mem

import (
	"sync"

	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/internal/types"
)

type lockServiceImpl struct {
	types.NwServiceImpl
	config common.Config
	sync.Mutex
}

func NewLockServiceImpl(nw types.Network, config common.Config) types.LockService {
	lockService := lockServiceImpl{
		NwServiceImpl: types.NwServiceImpl{
			Nw: nw,
		},
		config: config,
	}
	return &lockService
}

func (l *lockServiceImpl) Init() {

}
