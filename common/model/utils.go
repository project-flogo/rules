package model

//IdentifiersToString Take a slice of Identifiers and return a string representation
func IdentifiersToString(identifiers []TupleType) string {
	str := ""
	for _, idr := range identifiers {
		str += string(idr) + ", "
	}
	return str
}

// Contains returns true if an identifier exists in the identifier array
func Contains(identifiers []TupleType, toCheck TupleType) (bool, int) {
	for idx, id := range identifiers {
		if id == toCheck {
			return true, idx
		}
	}
	return false, -1
}
