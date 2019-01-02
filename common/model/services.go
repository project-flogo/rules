package model


type Service interface {
	Init()
}

type TupleStore interface {
	Service
	GetTupleByKey(key TupleKey) Tuple
	SaveTuple(tuple Tuple)
	SaveTuples(added map[string]map[string]Tuple)
	SaveModifiedTuples(modified map[string]map[string]RtcModified)
	DeleteTupleByStringKey(key TupleKey)
	DeleteTuples(deleted map[string]map[string]Tuple)
}
