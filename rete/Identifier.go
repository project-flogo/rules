package rete

//identifier An internal representation of a 'DataSource'
type identifier interface {
	getName() string
	equals(identifier) bool
	String() string
}

type identifierImpl struct {
	name string
	// alias string
}

func newIdentifier(name string /* alias string */) identifier {
	idr := identifierImpl{}
	idr.initIdentifierImpl(name /*alias */)
	return &idr
}
func (idr *identifierImpl) initIdentifierImpl(name string /*alias string*/) {
	idr.name = name
	// idr.alias = alias
}

func (idr *identifierImpl) getName() string {
	return idr.name
}

func (idr *identifierImpl) equals(other identifier) bool {
	return idr.name == other.getName() /*&& idr.alias == other.GetAlias()*/
}

func (idr *identifierImpl) String() string {
	// return "[" + idr.name + "," + idr.alias + "]"
	return "[" + idr.name + "]"
}
