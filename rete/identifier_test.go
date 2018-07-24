package rete

import (
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/bego/common/model"
)

func TestIdentifier(t *testing.T) {

	i1 := model.TupleTypeAlias("1")

	i2 := model.TupleTypeAlias("1")

	if i1 == i2 {
		fmt.Printf("yes they are equal!")
	} else {
		fmt.Printf("yes they are NOT equal!")
	}

	ids := []model.TupleTypeAlias{i1, i2}

	node := newNode(ids)

	ids2 := node.getIdentifiers()

	for _, n := range ids2 {
		fmt.Printf("%s", n)
	}

}
