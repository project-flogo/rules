package redis

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/redisutils"
)

type storeImpl struct {
	//allTuples map[string]model.Tuple
	prefix string
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

	strKey := ms.prefix + key.GetTupleDescriptor().Name + ":" + key.String()
	vals := hdl.HGetAll(strKey)


	tuple, err := model.NewTuple(model.TupleType(key.GetTupleDescriptor().Name), vals)

	if err != nil {
		return tuple
	}
	return nil
}

func (ms *storeImpl) SaveTuple(tuple model.Tuple) {
	m := tuple.ToMap()

	strKey := ms.prefix + tuple.GetTupleDescriptor().Name + ":" + tuple.GetKey().String()

	hdl := redisutils.GetRedisHdl()
	hdl.HSetAll(strKey, m)
}

func (ms *storeImpl) DeleteTupleByStringKey(key model.TupleKey) {
	strKey := ms.prefix + key.GetTupleDescriptor().Name + ":" + key.String()
	hdl := redisutils.GetRedisHdl()
	hdl.Del(strKey)
}

func (ms *storeImpl) SaveTuples(added map[string]map[string]model.Tuple) {

}

func (ms *storeImpl) SaveModifiedTuples(modified map[string]map[string]model.RtcModified) {

}

func (ms *storeImpl) DeleteTuples(deleted map[string]map[string]model.Tuple) {

}

