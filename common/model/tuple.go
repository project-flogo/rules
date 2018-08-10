package model

import (
	"context"
	"time"
)

var reteCTXKEY = RetecontextKeyType{}

//Tuple is a runtime representation of a data tuple
type Tuple interface {
	GetTypeAlias() TupleType
	GetString(name string) string
	GetInt(name string) int
	GetFloat(name string) float64
	GetDateTime(name string) time.Time
	GetProperties() []string
	GetTupleDescriptor() *TupleDescriptor
}

//MutableTuple mutable part of the tuple
type MutableTuple interface {
	Tuple
	SetString(ctx context.Context, name string, value string)
	SetInt(ctx context.Context, name string, value int)
	SetFloat(ctx context.Context, name string, value float64)
	SetDatetime(ctx context.Context, name string, value time.Time)
}

type tupleImpl struct {
	tupleType TupleType
	tuples    map[string]interface{}
	key       TupleKey
	td        *TupleDescriptor
}

func NewTuple(tupleType TupleType) MutableTuple {
	td := GetTupleDescriptor(tupleType)
	if td == nil {
		return nil
	}
	st := tupleImpl{}
	st.initTuple(td)
	st.key = newTupleKey(TupleType(td.Name))
	return &st
}

func (st *tupleImpl) initTuple(td *TupleDescriptor) {
	st.tuples = make(map[string]interface{})
	st.tupleType = TupleType(td.Name)
	st.td = td
}

func (st *tupleImpl) GetTypeAlias() TupleType {
	return st.tupleType
}

func (st *tupleImpl) GetString(name string) string {
	v := st.tuples[name]
	return v.(string)
}
func (st *tupleImpl) GetInt(name string) int {
	v := st.tuples[name]
	return v.(int)
}
func (st *tupleImpl) GetFloat(name string) float64 {
	v := st.tuples[name]
	return v.(float64)
}
func (st *tupleImpl) GetDateTime(name string) time.Time {
	v := st.tuples[name]
	return v.(time.Time)
}

func (st *tupleImpl) SetString(ctx context.Context, name string, value string) {
	if !st.validateUpdate(name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}
}

func (st *tupleImpl) SetInt(ctx context.Context, name string, value int) {
	if !st.validateUpdate(name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}
}
func (st *tupleImpl) SetFloat(ctx context.Context, name string, value float64) {
	if !st.validateUpdate(name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}
}
func (st *tupleImpl) SetDatetime(ctx context.Context, name string, value time.Time) {
	if !st.validateUpdate(name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}
}

func (st *tupleImpl) GetTupleDescriptor() *TupleDescriptor {
	return st.td
}

func callChangeListener(ctx context.Context, tuple Tuple, prop string) {
	if ctx != nil {
		ctxR := ctx.Value(reteCTXKEY)
		if ctxR != nil {
			valChangeLister := ctxR.(ValueChangeListener)
			valChangeLister.OnValueChange(tuple, prop)
		}
	}
}

func (st *tupleImpl) GetProperties() []string {
	keys := []string{}
	for k := range st.tuples {
		keys = append(keys, k)
	}
	return keys
}

func (st *tupleImpl) GetKey() TupleKey {
	keyImpl := st.key.(tupleKeyImpl)
	if keyImpl.keys != nil {
		return keyImpl
	} else {
		keyImpl.keys = make(map[string]interface{})
		for _, keyProp := range st.GetTupleDescriptor().GetKeyProps() {
			keyImpl.keys[keyProp] = st.tuples[keyProp]
		}
	}
	return keyImpl
}

func (st *tupleImpl) validateUpdate(name string, value interface{}) bool {
	//TODO: Check property's type and value's type compatibility
	_, ok := st.td.GetProperty(name)
	return ok
}
