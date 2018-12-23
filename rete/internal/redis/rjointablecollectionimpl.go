package redis

import "github.com/project-flogo/rules/rete/internal/types"

type joinTableCollectionImpl struct {
	allJoinTables map[int]types.JoinTable
}

func NewJoinTableCollection(config map[string]interface{}) types.JtService {
	jtc := joinTableCollectionImpl{}
	jtc.allJoinTables = make(map[int]types.JoinTable)
	return &jtc
}
func (jtc *joinTableCollectionImpl) Init() {

}

func (jtc *joinTableCollectionImpl) GetJoinTable(joinTableID int) types.JoinTable {
	return jtc.allJoinTables[joinTableID]
}

func (jtc *joinTableCollectionImpl) AddJoinTable(joinTable types.JoinTable) {
	jtc.allJoinTables[joinTable.GetID()] = joinTable
}

func (jtc *joinTableCollectionImpl) RemoveJoinTable(joinTable types.JoinTable) {
	delete(jtc.allJoinTables,joinTable.GetID())
}