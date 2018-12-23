package mem

import (
	"github.com/project-flogo/rules/rete/internal/types"
	"sync/atomic"
)

type idGenImpl struct {
	config    map[string]interface{}
	currentId int32
}

func NewIdImpl(config map[string]interface{}) types.IdGen {
	idg := idGenImpl{}
	idg.config = config
	idg.currentId = 0
	return &idg
}

func (id *idGenImpl) Init() {
	id.currentId = int32(id.GetMaxID())
}

func (id *idGenImpl) GetNextID() int {
	i := atomic.AddInt32(&id.currentId, 1)
	return int(i)
}

func (id *idGenImpl) GetMaxID() int {
	i := atomic.LoadInt32(&id.currentId)
	return int(i)
}
