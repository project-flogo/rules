package redis

import (
	"github.com/project-flogo/rules/rete/internal/types"
	"github.com/project-flogo/rules/common/model"
)

type joinTableCollectionImpl struct {
	allJoinTables map[string]types.JoinTable
}

func NewJoinTableCollection(config map[string]interface{}) types.JtService {
	jtc := joinTableCollectionImpl{}
	jtc.allJoinTables = make(map[string]types.JoinTable)
	return &jtc
}
func (jtc *joinTableCollectionImpl) Init() {

}

func (jtc *joinTableCollectionImpl) GetJoinTable(joinTableName string) types.JoinTable {
	return jtc.allJoinTables[joinTableName]
}

func (jtc *joinTableCollectionImpl) GetOrCreateJoinTable(nw types.Network, rule model.Rule, conditionVar model.Condition, identifiers []model.TupleType) types.JoinTable {
	jT, found := jtc.allJoinTables[conditionVar.GetName()]
	if !found {
		jTn := joinTableImpl{}
		jTn.initJoinTableImpl(nw, rule, identifiers)
		jtc.allJoinTables[conditionVar.GetName()] = &jTn
		jT = &jTn
	}
	return jT}

func (jtc *joinTableCollectionImpl) AddJoinTable(joinTable types.JoinTable) {
	jtc.allJoinTables[joinTable.GetName()] = joinTable
}

func (jtc *joinTableCollectionImpl) RemoveJoinTable(joinTable types.JoinTable) {
	delete(jtc.allJoinTables,joinTable.GetName())
}