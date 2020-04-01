package mem

import (
	"sync/atomic"

	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/internal/types"
)

type idGenServiceImpl struct {
	types.NwServiceImpl
	config    common.Config
	currentId int32
}

func NewIdGenImpl(nw types.Network, config common.Config) types.IdGen {
	idg := idGenServiceImpl{}
	idg.config = config
	idg.currentId = 0
	idg.Nw = nw
	return &idg
}

func (id *idGenServiceImpl) Init() {
	id.currentId = int32(id.GetMaxID())
}

func (id *idGenServiceImpl) GetNextID() int {
	i := atomic.AddInt32(&id.currentId, 1)
	return int(i)
}

func (id *idGenServiceImpl) GetMaxID() int {
	i := atomic.LoadInt32(&id.currentId)
	return int(i)
}
