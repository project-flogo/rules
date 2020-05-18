package rete

import "github.com/project-flogo/rules/common/model"

//These set operations are used in building the rete network. See Network.buildNetwork

//AppendIdentifiers ... Append identifiers from set2 to set1
func AppendIdentifiers(set1 []model.TupleType, set2 []model.TupleType) []model.TupleType {
	union := []model.TupleType{}
	union = append(union, set1...)
	union = append(union, set2...)
	return union
}

//ContainedByFirst ... true if second is a subset of first
func ContainedByFirst(first []model.TupleType, second []model.TupleType) bool {

	if len(second) == 0 {
		return true
	} else if len(first) == 0 {
		return false
	}
	for _, idFromSecond := range second {
		contains := false
		for _, idFromFirst := range first {
			if idFromSecond == idFromFirst {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}
	return true
}

//UnionOfOtherTwoContainsAllFromFirst ... true if union of second and third cover all of first
func UnionOfOtherTwoContainsAllFromFirst(first []model.TupleType, second []model.TupleType, third []model.TupleType) bool {
	return ContainedByFirst(UnionIdentifiers(second, third), first)
}

//UnionIdentifiers ... union of the first and second sets
func UnionIdentifiers(first []model.TupleType, second []model.TupleType) []model.TupleType {
	union := []model.TupleType{}
	union = append(union, first...)
	union = append(union, SecondMinusFirst(first, second)...)
	return union
}

//SecondMinusFirst ... returns elements in the second that arent in the first
func SecondMinusFirst(first []model.TupleType, second []model.TupleType) []model.TupleType {
	minus := []model.TupleType{}
outer:
	for _, idrSecond := range second {
		for _, idrFirst := range first {
			if idrSecond == idrFirst {
				continue outer
			}
		}
		minus = append(minus, idrSecond)
	}
	return minus
}

//IntersectionIdentifiers .. intersection of the two sets
func IntersectionIdentifiers(first []model.TupleType, second []model.TupleType) []model.TupleType {
	intersect := []model.TupleType{}
	for _, idrSecond := range second {
		for _, idrFirst := range first {
			if idrSecond == idrFirst {
				intersect = append(intersect, idrSecond)
			}
		}
	}
	return intersect
}

//EqualSets ... compare two identifiers based on their contents
func EqualSets(first []model.TupleType, second []model.TupleType) bool {
	return len(SecondMinusFirst(first, second)) == 0 &&
		len(SecondMinusFirst(second, first)) == 0
}

//GetIndex ... return the index of thisIdr in identifiers
func GetIndex(identifiers []model.TupleType, thisIdr model.TupleType) int {
	for i, idr := range identifiers {
		if idr == thisIdr {
			return i
		}
	}
	return -1
}
