package model

import (
	"context"
	"time"
)

type tupleImpl struct {
	tupleType TupleType
	tuples    map[string]interface{}
}

func NewTuple(tuple TupleType) MutableStreamTuple {
	st := tupleImpl{}
	st.initStreamTuple(tuple)
	return &st
}

func (st *tupleImpl) initStreamTuple(tupleType TupleType) {
	st.tuples = make(map[string]interface{})
	st.tupleType = tupleType
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

func (st *tupleImpl) SetString(ctx context.Context, rs RuleSession, name string, value string) {
	if rs == nil || rs != nil && !rs.ValidateUpdate(st.tupleType, name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}

}
func (st *tupleImpl) SetInt(ctx context.Context, rs RuleSession, name string, value int) {
	if rs == nil || rs != nil && !rs.ValidateUpdate(st.tupleType, name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}
}
func (st *tupleImpl) SetFloat(ctx context.Context, rs RuleSession, name string, value float64) {
	if rs == nil || rs != nil && !rs.ValidateUpdate(st.tupleType, name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}
}
func (st *tupleImpl) SetDatetime(ctx context.Context, rs RuleSession, name string, value time.Time) {
	if rs == nil || rs != nil && !rs.ValidateUpdate(st.tupleType, name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}
}

func callChangeListener(ctx context.Context, tuple Tuple, prop string) {
	if ctx != nil {
		ctxR := ctx.Value(reteCTXKEY)
		if ctxR != nil {
			valChangeLister := ctxR.(ValueChangeHandler)
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