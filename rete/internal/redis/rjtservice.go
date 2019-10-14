package redis

import (
	"sync"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/internal/types"
)

type jtServiceImpl struct {
	types.NwServiceImpl
	allJoinTables map[string]types.JoinTable
	sync.RWMutex
}

func NewJoinTableCollection(nw types.Network, config common.Config) types.JtService {
	jtc := jtServiceImpl{}
	jtc.Nw = nw
	jtc.allJoinTables = make(map[string]types.JoinTable)
	return &jtc
}
func (jtc *jtServiceImpl) Init() {

}

func (jtc *jtServiceImpl) GetJoinTable(name string) types.JoinTable {
	jtc.RLock()
	defer jtc.RUnlock()
	return jtc.allJoinTables[name]
}

func (jtc *jtServiceImpl) GetOrCreateJoinTable(nw types.Network, rule model.Rule, identifiers []model.TupleType, name string) types.JoinTable {
	jtc.Lock()
	defer jtc.Unlock()
	jT, found := jtc.allJoinTables[name]
	if !found {
		jT = newJoinTableImpl(nw, rule, identifiers, name)
		jtc.allJoinTables[name] = jT
	}
	return jT
}
