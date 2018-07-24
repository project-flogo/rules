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

func copyIntoTupleMap(handles []reteHandle) map[model.TupleTypeAlias]model.StreamTuple {
	tupleMap := map[model.TupleTypeAlias]model.StreamTuple{}
	tuples := make([]model.StreamTuple, len(handles))
	for i := 0; i < len(handles); i++ {
		tuples[i] = handles[i].getTuple()
		tupleMap[tuples[i].GetTypeAlias()] = tuples[i] //assuming no self-joins! need to correct this!
	}
	return tupleMap
}

func convertToTupleMap(tuples []model.StreamTuple) map[model.TupleTypeAlias]model.StreamTuple {
	tupleMap := map[model.TupleTypeAlias]model.StreamTuple{}
	for i := 0; i < len(tuples); i++ {
		tupleMap[tuples[i].GetTypeAlias()] = tuples[i]
	}
	return tupleMap
}
