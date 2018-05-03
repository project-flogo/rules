package rete

import (
	"testing"
)

func TestIdentifierUtils(t *testing.T) {

	first(t)
	second(t)
	third(t)
	fourth(t)

}
func first(t *testing.T) {
	first := []identifier{newIdentifier("1"), newIdentifier("2")}
	second := []identifier{newIdentifier("1"), newIdentifier("2")}
	if len(UnionIdentifiers(first, second)) != 2 {
		t.Error("Failed")
	}
	if len(SecondMinusFirst(first, second)) != 0 {
		t.Error("Failed")
	}
	if len(IntersectionIdentifiers(first, second)) != 2 {
		t.Error("Failed")
	}
	if !ContainedByFirst(first, second) {
		t.Error("Failed")
	}

	if GetIndex(first, newIdentifier("1")) != 0 {
		t.Error("Failed")
	}
	if GetIndex(first, newIdentifier("2")) != 1 {
		t.Error("Failed")
	}
}

func second(t *testing.T) {
	first := []identifier{newIdentifier("1")}
	second := []identifier{newIdentifier("1"), newIdentifier("2")}
	if len(UnionIdentifiers(first, second)) != 2 {
		t.Error("Failed")
	}
	if len(SecondMinusFirst(first, second)) != 1 {
		t.Error("Failed")
	}
	if len(IntersectionIdentifiers(first, second)) != 1 {
		t.Error("Failed")
	}
	if ContainedByFirst(first, second) {
		t.Error("Failed")
	}
}

func third(t *testing.T) {
	first := []identifier{newIdentifier("1"), newIdentifier("2")}
	second := []identifier{newIdentifier("1")}
	if len(UnionIdentifiers(first, second)) != 2 {
		t.Error("Failed")
	}
	if len(SecondMinusFirst(first, second)) != 0 {
		t.Error("Failed")
	}
	if len(IntersectionIdentifiers(first, second)) != 1 {
		t.Error("Failed")
	}
	if !ContainedByFirst(first, second) {
		t.Error("Failed")
	}
}

func fourth(t *testing.T) {
	first := []identifier{newIdentifier("1"), newIdentifier("2")}
	second := []identifier{newIdentifier("1"), newIdentifier("2")}
	third := []identifier{newIdentifier("1"), newIdentifier("2"),
		newIdentifier("3")}

	if OtherTwoAreContainedByFirst(first, second, third) {
		t.Error("Failed")
	}

	if !OtherTwoAreContainedByFirst(third, second, first) {
		t.Error("Failed")
	}
}
