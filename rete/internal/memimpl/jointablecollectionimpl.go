package memimpl

import "github.com/project-flogo/rules/rete/internal/types"

type joinTableCollectionImpl struct {
	allJoinTables map[int]types.JoinTable
}
