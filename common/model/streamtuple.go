package model

import (
	"context"
	"time"
)

//StreamTuple is a runtime representation of a stream of data
type StreamTuple interface {
	GetStreamDataSource() StreamSource
	GetString(name string) string
	GetInt(name string) int
	GetFloat(name string) float64
	GetDateTime(name string) time.Time
}

//MutableStreamTuple mutable part of the stream tuple
type MutableStreamTuple interface {
	StreamTuple
	SetString(ctx context.Context, name string, value string)
	SetInt(ctx context.Context, name string, value int)
	SetFloat(ctx context.Context, name string, value float64)
	SetDatetime(ctx context.Context, name string, value time.Time)
}

type streamImpl struct {
	dataSource StreamSource
	tuples     map[string]interface{}
}

func NewStreamTuple(dataSource StreamSource) MutableStreamTuple {
	streamImplVar := streamImpl{}
	streamImplVar.initStreamTuple(dataSource)
	return &streamImplVar
}

func (streamImplVar *streamImpl) initStreamTuple(dataSource StreamSource) {
	streamImplVar.tuples = make(map[string]interface{})
	streamImplVar.dataSource = dataSource
}

func (streamImplVar *streamImpl) GetStreamDataSource() StreamSource {
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
	streamImplVar.tuples[name] = value
	callChangeListener(ctx, streamImplVar)
}
func (streamImplVar *streamImpl) SetInt(ctx context.Context, name string, value int) {
	streamImplVar.tuples[name] = value
	callChangeListener(ctx, streamImplVar)
}
func (streamImplVar *streamImpl) SetFloat(ctx context.Context, name string, value float64) {
	streamImplVar.tuples[name] = value
	callChangeListener(ctx, streamImplVar)
}
func (streamImplVar *streamImpl) SetDatetime(ctx context.Context, name string, value time.Time) {
	streamImplVar.tuples[name] = value
	callChangeListener(ctx, streamImplVar)
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
