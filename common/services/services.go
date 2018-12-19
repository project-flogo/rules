package services

import "github.com/project-flogo/rules/common/model"

type TupleStore interface {
	GetTupleByStringKey(key string) model.Tuple
	SaveTuple(tuple model.Tuple)
	DeleteTupleByStringKey(key string)
}
