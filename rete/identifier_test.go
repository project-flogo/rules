package rete

import (
	"fmt"
	"testing"

	"github.com/project-flogo/rules/common/model"
)

func TestIdentifier(t *testing.T) {

	i1 := model.TupleType("1")

	i2 := model.TupleType("1")

	if i1 == i2 {
		fmt.Printf("yes they are equal!")
	} else {
		fmt.Printf("yes they are NOT equal!")
	}

	ids := []model.TupleType{i1, i2}

	node := newNode(nil, nil, ids)

	ids2 := node.getIdentifiers()

	for _, n := range ids2 {
		fmt.Printf("%s", n)
	}

}
