package model

import (
	"context"
	"time"
)

type streamImpl struct {
	dataSource TupleTypeAlias
	tuples     map[string]interface{}
}

func NewStreamTuple(dataSource TupleTypeAlias) MutableStreamTuple {
	streamImplVar := streamImpl{}
	streamImplVar.initStreamTuple(dataSource)
	return &streamImplVar
}

func (streamImplVar *streamImpl) initStreamTuple(dataSource TupleTypeAlias) {
	streamImplVar.tuples = make(map[string]interface{})
	streamImplVar.dataSource = dataSource
}

func (streamImplVar *streamImpl) GetTypeAlias() TupleTypeAlias {
	return streamImplVar.dataSource
}

func (streamImplVar *streamImpl) GetString(name string) string {
	v := streamImplVar.tuples[name]
	return v.(string)
}
func (streamImplVar *streamImpl) GetInt(name string) int {
	v := streamImplVar.tuples[name]
	return v.(int)
}
func (streamImplVar *streamImpl) GetFloat(name string) float64 {
	v := streamImplVar.tuples[name]
	return v.(float64)
}
func (streamImplVar *streamImpl) GetDateTime(name string) time.Time {
	v := streamImplVar.tuples[name]
	return v.(time.Time)
}

func (streamImplVar *streamImpl) SetString(ctx context.Context, name string, value string) {
	if streamImplVar.tuples[name] != value {
		streamImplVar.tuples[name] = value
		callChangeListener(ctx, streamImplVar)
	}
}
func (streamImplVar *streamImpl) SetInt(ctx context.Context, name string, value int) {
	if streamImplVar.tuples[name] != value {
		streamImplVar.tuples[name] = value
		callChangeListener(ctx, streamImplVar)
	}
}
func (streamImplVar *streamImpl) SetFloat(ctx context.Context, name string, value float64) {
	if streamImplVar.tuples[name] != value {
		streamImplVar.tuples[name] = value
		callChangeListener(ctx, streamImplVar)
	}
}
func (streamImplVar *streamImpl) SetDatetime(ctx context.Context, name string, value time.Time) {
	if streamImplVar.tuples[name] != value {
		streamImplVar.tuples[name] = value
		callChangeListener(ctx, streamImplVar)
	}
}

func callChangeListener(ctx context.Context, tuple StreamTuple) {
	if ctx != nil {
		valChangeListerI := ctx.Value(reteCTXKEY)
		if valChangeListerI != nil {
			valChangeLister := valChangeListerI.(ValueChangeHandler)
			valChangeLister.OnValueChange(tuple)
		}
	}
}
