package rete

//translation utilities between handles/tuples to pass to user conditions and actions

import "github.com/TIBCOSoftware/bego/common/model"

func copyIntoTupleArray(handles []reteHandle) []model.StreamTuple {
	tuples := make([]model.StreamTuple, len(handles))
	for i := 0; i < len(handles); i++ {
		tuples[i] = handles[i].getTuple()
	}
	return tuples
}

func copyIntoTupleMap(handles []reteHandle) map[model.StreamSource]model.StreamTuple {
	tupleMap := map[model.StreamSource]model.StreamTuple{}
	tuples := make([]model.StreamTuple, len(handles))
	for i := 0; i < len(handles); i++ {
		tuples[i] = handles[i].getTuple()
		tupleMap[tuples[i].GetStreamDataSource()] = tuples[i]
	}
	return tupleMap
}

func convertToTupleMap(tuples []model.StreamTuple) map[model.StreamSource]model.StreamTuple {
	tupleMap := map[model.StreamSource]model.StreamTuple{}
	for i := 0; i < len(tuples); i++ {
		tupleMap[tuples[i].GetStreamDataSource()] = tuples[i]
	}
	return tupleMap
}
