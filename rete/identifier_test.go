package rete

import (
	"fmt"
	"testing"
)

func TestIdentifier(t *testing.T) {

	i1 := newIdentifier("1")

	i2 := newIdentifier("1")

	if i1.equals(i2) {
		fmt.Printf("yes they are equal!")
	} else {
		fmt.Printf("yes they are NOT equal!")
	}

	ids := []identifier{i1, i2}

	node := newNode(ids)

	ids2 := node.getIdentifiers()

	for _, n := range ids2 {
		fmt.Printf("%s", n)
	}

}
