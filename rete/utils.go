package rete

//translation utilities between handles/tuples to pass to user conditions and actions

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
)

func copyIntoTupleArray(handles []types.ReteHandle) []model.Tuple {
	tuples := make([]model.Tuple, len(handles))
	for i := 0; i < len(handles); i++ {
		tuples[i] = handles[i].GetTuple()
	}
	return tuples
}

func copyIntoTupleMap(handles []types.ReteHandle) map[model.TupleType]model.Tuple {
	tupleMap := map[model.TupleType]model.Tuple{}
	tuples := make([]model.Tuple, len(handles))
	for i := 0; i < len(handles); i++ {
		tuples[i] = handles[i].GetTuple()
		tupleMap[tuples[i].GetTupleType()] = tuples[i] //assuming no self-joins! need to correct this!
	}
	return tupleMap
}

func convertToTupleMap(tuples []model.Tuple) map[model.TupleType]model.Tuple {
	tupleMap := map[model.TupleType]model.Tuple{}
	for i := 0; i < len(tuples); i++ {
		tupleMap[tuples[i].GetTupleType()] = tuples[i]
	}
	return tupleMap
}
