package mem

import (
	"fmt"
	"sync"

	"github.com/project-flogo/rules/common/model"
)

type storeImpl struct {
	allTuples map[string]model.Tuple
	sync.RWMutex
}

func NewStore(jsonConfig map[string]interface{}) model.TupleStore {
	ms := storeImpl{
		allTuples: make(map[string]model.Tuple),
	}
	return &ms
}

func (ms *storeImpl) Init() {

}

func (ms *storeImpl) GetTupleByKey(key model.TupleKey) model.Tuple {
	ms.RLock()
	defer ms.RUnlock()
	return ms.allTuples[key.String()]
}

func (ms *storeImpl) SaveTuple(tuple model.Tuple) {
	ms.Lock()
	defer ms.Unlock()
	ms.allTuples[tuple.GetKey().String()] = tuple
}

func (ms *storeImpl) DeleteTuple(key model.TupleKey) {
	ms.Lock()
	defer ms.Unlock()
	delete(ms.allTuples, key.String())
}

func (ms *storeImpl) SaveTuples(added map[string]map[string]model.Tuple) {
	for tupleType, tuples := range added {
		for key, tuple := range tuples {
			fmt.Printf("Saving tuple. Type [%s] Key [%s], Val [%v]\n", tupleType, key, tuple)
			ms.SaveTuple(tuple)
		}
	}
}

func (ms *storeImpl) SaveModifiedTuples(modified map[string]map[string]model.RtcModified) {
	for tupleType, mmap := range modified {
		for key, mdfd := range mmap {
			fmt.Printf("Saving tuple. Type [%s] Key [%s], Val [%v]\n", tupleType, key, mdfd.GetTuple())
			ms.SaveTuple(mdfd.GetTuple())
		}
	}
}

func (ms *storeImpl) DeleteTuples(deleted map[string]map[string]model.Tuple) {
	for tupleType, tuples := range deleted {
		for key, tuple := range tuples {
			fmt.Printf("Deleting tuple. Type [%s] Key [%s], Val [%v]\n", tupleType, key, tuple)
			ms.DeleteTuple(tuple.GetKey())
		}
	}
}
