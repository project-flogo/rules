package rete

//These set operations are used in building the rete network. See Network.buildNetwork

//AppendIdentifiers ... Append identifiers from set2 to set1
func AppendIdentifiers(set1 []identifier, set2 []identifier) []identifier {
	union := []identifier{}
	union = append(union, set1...)
	union = append(union, set2...)
	return union
}

//ContainedByFirst ... true if second is a subset of first
func ContainedByFirst(first []identifier, second []identifier) bool {

	if len(second) == 0 {
		return true
	} else if len(first) == 0 {
		return false
	}
	for _, idFromSecond := range second {
		contains := false
		for _, idFromFirst := range first {
			if idFromSecond.equals(idFromFirst) {
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

//OtherTwoAreContainedByFirst ... true if second and third are part of first
func OtherTwoAreContainedByFirst(first []identifier, second []identifier, third []identifier) bool {
	return ContainedByFirst(first, second) && ContainedByFirst(first, third)
}

//UnionIdentifiers ... union of the first and second sets
func UnionIdentifiers(first []identifier, second []identifier) []identifier {
	union := []identifier{}
	union = append(union, first...)
	union = append(union, SecondMinusFirst(first, second)...)
	return union
}

//SecondMinusFirst ... returns elements in the second that arent in the first
func SecondMinusFirst(first []identifier, second []identifier) []identifier {
	minus := []identifier{}
outer:
	for _, idrSecond := range second {
		for _, idrFirst := range first {
			if idrSecond.equals(idrFirst) {
				continue outer
			}
		}
		minus = append(minus, idrSecond)
	}
	return minus
}

//IntersectionIdentifiers .. intersection of the two sets
func IntersectionIdentifiers(first []identifier, second []identifier) []identifier {
	intersect := []identifier{}
	for _, idrSecond := range second {
		for _, idrFirst := range first {
			if idrSecond.equals(idrFirst) {
				intersect = append(intersect, idrSecond)
			}
		}
	}
	return intersect
}

//EqualSets ... compare two identifiers based on their contents
func EqualSets(first []identifier, second []identifier) bool {
	return len(SecondMinusFirst(first, second)) == 0 && len(SecondMinusFirst(first, second)) == 0
}

//GetIndex ... return the index of thisIdr in identifiers
func GetIndex(identifiers []identifier, thisIdr identifier) int {
	for i, idr := range identifiers {
		if idr.equals(thisIdr) {
			return i
		}
		i++
	}
	return -1
}

//IdentifiersToString Take a slice of Identifiers and return a string representation
func IdentifiersToString(identifiers []identifier) string {
	str := ""
	for _, idr := range identifiers {
		str += idr.String() + ", "
	}
	return str
}
