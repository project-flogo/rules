package redis

import (
	"fmt"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/common"
)

type storeImpl struct {
	//allTuples map[string]model.Tuple
	prefix string
	config common.Config
	redisutils.RedisHdl
}

func NewStore(config common.Config) model.TupleStore {
	ms := storeImpl{
		config:   config,
		RedisHdl: redisutils.NewRedisHdl(config.Stores.Redis),
	}
	return &ms
}

func (ms *storeImpl) Init() {

}

func (ms *storeImpl) GetTupleByKey(key model.TupleKey) model.Tuple {
	strKey := ms.prefix + key.String()
	vals := ms.HGetAll(strKey)

	tuple, err := model.NewTuple(model.TupleType(key.GetTupleDescriptor().Name), vals)

	if err == nil {
		return tuple
	}
	return nil
}

func (ms *storeImpl) SaveTuple(tuple model.Tuple) {
	m := tuple.ToMap()

	strKey := ms.prefix + tuple.GetKey().String()

	ms.HSetAll(strKey, m)
}

func (ms *storeImpl) DeleteTuple(key model.TupleKey) {
	strKey := ms.prefix + key.String()
	ms.Del(strKey)
}

func (ms *storeImpl) SaveTuples(added map[string]map[string]model.Tuple) {
	for tupleType, tuples := range added {
		for key, tuple := range tuples {
			fmt.Printf("Saving tuple. Type [%s] Key [%s]\n", tupleType, key)
			strKey := ms.prefix + key
			ms.HSetAll(strKey, tuple.ToMap())
		}
	}
}

func (ms *storeImpl) SaveModifiedTuples(modified map[string]map[string]model.RtcModified) {
	for tupleType, mmap := range modified {
		for key, mdfd := range mmap {
			fmt.Printf("Saving tuple. Type [%s] Key [%s]\n", tupleType, key)
			strKey := ms.prefix + key
			ms.HSetAll(strKey, mdfd.GetTuple().ToMap())
		}
	}
}

func (ms *storeImpl) DeleteTuples(deleted map[string]map[string]model.Tuple) {
	for tupleType, tuples := range deleted {
		for key, _ := range tuples {
			fmt.Printf("Deleting tuple. Type [%s] Key [%s]\n", tupleType, key)
			strKey := ms.prefix + key
			ms.Del(strKey)
		}
	}
}
