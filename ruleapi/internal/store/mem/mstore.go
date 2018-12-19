package mem

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/common/services"
)

type store struct {
	allTuples map[string]model.Tuple
}

func NewStore() services.TupleStore {
	ms := store{}
	ms.allTuples = make(map[string]model.Tuple)
	return &ms
}

func (ms *store) GetTupleByStringKey(key string) model.Tuple {
	return ms.allTuples[key]
}
func (ms *store) SaveTuple(tuple model.Tuple) {
	ms.allTuples[tuple.GetKey().String()] = tuple
}
func (ms *store) DeleteTupleByStringKey(key string) {
	delete(ms.allTuples, key)
}
