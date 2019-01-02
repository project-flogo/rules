package mem

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
)

type jtServiceImpl struct {
	types.NwServiceImpl
	allJoinTables map[string]types.JoinTable
}

func NewJoinTableCollection(nw types.Network, config map[string]interface{}) types.JtService {
	jtc := jtServiceImpl{}
	jtc.Nw = nw
	jtc.allJoinTables = make(map[string]types.JoinTable)
	return &jtc
}
func (jtc *jtServiceImpl) Init() {

}

func (jtc *jtServiceImpl) GetJoinTable(joinTableName string) types.JoinTable {
	return jtc.allJoinTables[joinTableName]
}

func (jtc *jtServiceImpl) AddJoinTable(joinTable types.JoinTable) {
	jtc.allJoinTables[joinTable.GetName()] = joinTable
}

func (jtc *jtServiceImpl) RemoveJoinTable(jtName string) {
	delete(jtc.allJoinTables, jtName)
}

func (jtc *jtServiceImpl) GetOrCreateJoinTable(nw types.Network, rule model.Rule, identifiers []model.TupleType, name string) types.JoinTable {
	jT, found := jtc.allJoinTables[name]
	if !found {
		jT = newJoinTableImpl(nw, rule, identifiers, name)
		jtc.allJoinTables[name] = jT
	}
	return jT
}
