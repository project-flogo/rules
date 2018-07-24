package utils

//type mapImpl struct {
//	kvstore map[string]interface{}
//}
//
//func NewHashMap() Map {
//	mapImplVar := mapImpl{}
//	mapImplVar.initHashMap()
//	return &mapImplVar
//}
//
//func (mapImplVar *mapImpl) initHashMap() {
//	mapImplVar.kvstore = make(map[string]interface{})
//}
//
//func (mapImplVar *mapImpl) Put(key string, value interface{}) {
//	mapImplVar.kvstore[key] = value
//}
//
//func (mapImplVar *mapImpl) Get(key string) interface{} {
//	value := mapImplVar.kvstore[key]
//	return value
//}
//
//func (mapImplVar *mapImpl) Len() int {
//	return len(mapImplVar.kvstore)
//}
//
//func (mapImplVar *mapImpl) Remove(key string) interface{} {
//	value, ok := mapImplVar.kvstore[key]
//	if ok {
//		delete(mapImplVar.kvstore, key)
//	}
//	return value
//}
//
//func (mapImplVar *mapImpl) GetMap() map[string]interface{} {
//	return mapImplVar.kvstore
//}
//
//func (mapImplVar *mapImpl) ForEach(mapEntryHandler MapEntryHandler, context []interface{}) {
//	for k, v := range mapImplVar.GetMap() {
//		mapEntryHandler(k, v, context)
//	}
//}
