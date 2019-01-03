package redis

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
)

type jtServiceImpl struct {
	types.NwServiceImpl

	allJoinTables map[string]types.JoinTable
	prefix        string
}

func NewJoinTableCollection(nw types.Network, config map[string]interface{}) types.JtService {
	jtc := jtServiceImpl{}
	jtc.Nw = nw
	jtc.allJoinTables = make(map[string]types.JoinTable)
	jtc.prefix = nw.GetPrefix() + "jtc:"
	return &jtc
}
func (jtc *jtServiceImpl) Init() {

}

func (jtc *jtServiceImpl) GetJoinTable(name string) types.JoinTable {
	return jtc.allJoinTables[name]
}

func (jtc *jtServiceImpl) GetOrCreateJoinTable(nw types.Network, rule model.Rule, identifiers []model.TupleType, name string) types.JoinTable {
	jT, found := jtc.allJoinTables[name]
	if !found {
		jT = newJoinTableImpl(nw, rule, identifiers, name)
		jtc.allJoinTables[name] = jT
	}
	return jT
}

//func (jtc *jtServiceImpl) AddJoinTable(joinTable types.JoinTable) {
//	jtc.allJoinTables[joinTable.GetName()] = joinTable
//}
//
//func (jtc *jtServiceImpl) RemoveJoinTable(jtName string) {
//	delete(jtc.allJoinTables, jtName)
//}
