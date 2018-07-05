package utils

type lhMapImpl struct {
	mapImpl
	orderedKeyList ArrayList
}

func NewLinkedHashMap() Map {
	lhMapImplVar := lhMapImpl{}
	lhMapImplVar.initHashMap()
	lhMapImplVar.orderedKeyList = NewArrayList()
	return &lhMapImplVar
}

func (lhMapImplVar *lhMapImpl) Put(key string, value interface{}) {
	lhMapImplVar.mapImpl.Put(key, value)
	lhMapImplVar.orderedKeyList.Add(key)
}

func (lhMapImplVar *lhMapImpl) Get(key string) interface{} {
	return lhMapImplVar.mapImpl.Get(key)
}

func (lhMapImplVar *lhMapImpl) Len() int {
	return lhMapImplVar.mapImpl.Len()
}

func (lhMapImplVar *lhMapImpl) Remove(key string) interface{} {
	value := lhMapImplVar.mapImpl.Remove(key) //kvstore[key]
	lhMapImplVar.orderedKeyList.Remove(key)
	return value
}

func (lhMapImplVar *lhMapImpl) GetMap() map[string]interface{} {
	return lhMapImplVar.mapImpl.kvstore
}

func myListHandleEntry(entry interface{}, context []interface{}) {
	key := entry.(string)
	mapEntryHandler := context[0].(MapEntryHandler)
	lhMapImplVar := context[1].(*lhMapImpl)
	var origContext []interface{}
	if len(context) > 2 {
		origContext = context[2:]
	}
	value := lhMapImplVar.Get(key)
	if value != nil {
		mapEntryHandler(key, value, origContext)
	}
}

func (lhMapImplVar *lhMapImpl) ForEach(mapEntryHandler MapEntryHandler, context []interface{}) {
	context1 := make([]interface{}, 0)
	context1 = append(context1, mapEntryHandler) //[0]
	context1 = append(context1, lhMapImplVar)    //[1]
	if context != nil {
		context1 = append(context1, context...) //[context]
	}

	lhMapImplVar.orderedKeyList.ForEach(myListHandleEntry, context1)
}
