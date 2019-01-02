package mem

import (
	"github.com/project-flogo/rules/rete/internal/types"
	"sync/atomic"
)

type idGenServiceImpl struct {
	types.NwServiceImpl
	config    map[string]interface{}
	currentId int32
}

func NewIdGenImpl(nw types.Network, config map[string]interface{}) types.IdGen {
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
