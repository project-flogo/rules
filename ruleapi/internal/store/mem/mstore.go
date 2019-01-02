package mem

import (
	"github.com/project-flogo/rules/common/model"
)

type storeImpl struct {
	allTuples map[string]model.Tuple
}

func NewStore(jsonConfig map[string]interface{}) model.TupleStore {
	ms := storeImpl{}

	ms.allTuples = make(map[string]model.Tuple)
	return &ms
}

func (ms *storeImpl) Init() {

}

func (ms *storeImpl) GetTupleByKey(key model.TupleKey) model.Tuple {
	return ms.allTuples[key.String()]
}

func (ms *storeImpl) SaveTuple(tuple model.Tuple) {
	ms.allTuples[tuple.GetKey().String()] = tuple
}

func (ms *storeImpl) DeleteTupleByStringKey(key model.TupleKey) {
	delete(ms.allTuples, key.String())
}

func (ms *storeImpl) SaveTuples(added map[string]map[string]model.Tuple) {

}

func (ms *storeImpl) SaveModifiedTuples(modified map[string]map[string]model.RtcModified) {

}

func (ms *storeImpl) DeleteTuples(deleted map[string]map[string]model.Tuple) {

}
