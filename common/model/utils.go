package model

//IdentifiersToString Take a slice of Identifiers and return a string representation
func IdentifiersToString(identifiers []TupleTypeAlias) string {
	str := ""
	for _, idr := range identifiers {
		str += string(idr) + ", "
	}
	return str

}
