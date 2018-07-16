package rete

type retecontextKeyType struct {
}

var reteCTXKEY = retecontextKeyType{}

//store any context, may not know all keys upfront
type reteCtxValType struct {
	context map[interface{}]interface{}
}

func NewReteCtx() *reteCtxValType {
	reteCtxVal := reteCtxValType{}
	reteCtxVal.context = make(map[interface{}]interface{})
	return &reteCtxVal
}
