package redis

import (
	"fmt"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/redisutils"
)

type storeImpl struct {
	//allTuples map[string]model.Tuple
	prefix     string
	jsonConfig map[string]interface{}
}

func NewStore(jsonConfig map[string]interface{}) model.TupleStore {
	ms := storeImpl{}
	ms.jsonConfig = jsonConfig
	return &ms
}

func (ms *storeImpl) Init() {
	//ms.allTuples = make(map[string]model.Tuple)
	reteCfg := ms.jsonConfig["rs"].(map[string]interface{})
	ms.prefix = reteCfg["prefix"].(string)
	ms.prefix = ms.prefix + ":" + "s:"

}

func (ms *storeImpl) GetTupleByKey(key model.TupleKey) model.Tuple {
	hdl := redisutils.GetRedisHdl()

	strKey := ms.prefix + key.String()
	vals := hdl.HGetAll(strKey)

	tuple, err := model.NewTuple(model.TupleType(key.GetTupleDescriptor().Name), vals)

	if err == nil {
		return tuple
	}
	return nil
}

func (ms *storeImpl) SaveTuple(tuple model.Tuple) {
	m := tuple.ToMap()

	strKey := ms.prefix + tuple.GetKey().String()

	hdl := redisutils.GetRedisHdl()
	hdl.HSetAll(strKey, m)
}

func (ms *storeImpl) DeleteTuple(key model.TupleKey) {
	strKey := ms.prefix + key.String()
	hdl := redisutils.GetRedisHdl()
	hdl.Del(strKey)
}

func (ms *storeImpl) SaveTuples(added map[string]map[string]model.Tuple) {
	hdl := redisutils.GetRedisHdl()
	for tupleType, tuples := range added {
		for key, tuple := range tuples {
			fmt.Printf("Saving tuple. Type [%s] Key [%s]\n", tupleType, key)
			strKey := ms.prefix + key
			hdl.HSetAll(strKey, tuple.ToMap())
		}
	}
}

func (ms *storeImpl) SaveModifiedTuples(modified map[string]map[string]model.RtcModified) {
	hdl := redisutils.GetRedisHdl()
	for tupleType, mmap := range modified {
		for key, mdfd := range mmap {
			fmt.Printf("Saving tuple. Type [%s] Key [%s]\n", tupleType, key)
			strKey := ms.prefix + key
			hdl.HSetAll(strKey, mdfd.GetTuple().ToMap())
		}
	}
}

func (ms *storeImpl) DeleteTuples(deleted map[string]map[string]model.Tuple) {
	hdl := redisutils.GetRedisHdl()
	for tupleType, tuples := range deleted {
		for key, _ := range tuples {
			fmt.Printf("Deleting tuple. Type [%s] Key [%s]\n", tupleType, key)
			strKey := ms.prefix + key
			hdl.Del(strKey)
		}
	}
}
