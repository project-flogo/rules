package rete

import (
	"testing"

	"github.com/project-flogo/rules/common/model"
)

func TestIdentifierUtils(t *testing.T) {

	first(t)
	second(t)
	third(t)
	fourth(t)

}
func first(t *testing.T) {
	first := []model.TupleType{model.TupleType("1"), model.TupleType("2")}
	second := []model.TupleType{model.TupleType("1"), model.TupleType("2")}
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

	if GetIndex(first, model.TupleType("1")) != 0 {
		t.Error("Failed")
	}
	if GetIndex(first, model.TupleType("2")) != 1 {
		t.Error("Failed")
	}
}

func second(t *testing.T) {
	first := []model.TupleType{model.TupleType("1")}
	second := []model.TupleType{model.TupleType("1"), model.TupleType("2")}
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
	first := []model.TupleType{model.TupleType("1"), model.TupleType("2")}
	second := []model.TupleType{model.TupleType("1")}
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
	first := []model.TupleType{model.TupleType("1"), model.TupleType("2"), model.TupleType("3")}
	second := []model.TupleType{model.TupleType("1"), model.TupleType("2")}
	third := []model.TupleType{model.TupleType("3")}
	fourth := []model.TupleType{model.TupleType("2")}

	if !UnionOfOtherTwoContainsAllFromFirst(first, second, third) {
		t.Error("Failed")
	}

	if UnionOfOtherTwoContainsAllFromFirst(first, third, fourth) {
		t.Error("Failed")
	}
}
