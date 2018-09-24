package rete

import (
	"testing"
	"fmt"
)

func TestOne(t *testing.T) {
	fmt.Printf ("[is new] %d\n", rtcIsAsserted)
	fmt.Printf ("[is new] %d\n", rtcIsModified)
	fmt.Printf ("[is new] %d\n", rtcIsRetracted)

	i := 4 | rtcIsAsserted

	fmt.Printf ("[i is] %d\n", i)
}
