package services

import "github.com/project-flogo/rules/common/model"

type Service interface {
	Init()
}

type TupleStore interface {
	Service
	GetTupleByStringKey(key string) model.Tuple
	SaveTuple(tuple model.Tuple)
	DeleteTupleByStringKey(key string)
}
