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
	identifierImplVar := identifierImpl{}
	identifierImplVar.initIdentifierImpl(name /*alias */)
	return &identifierImplVar
}
func (identifierImplVar *identifierImpl) initIdentifierImpl(name string /*alias string*/) {
	identifierImplVar.name = name
	// identifierImplVar.alias = alias
}

func (identifierImplVar *identifierImpl) getName() string {
	return identifierImplVar.name
}

func (identifierImplVar *identifierImpl) equals(other identifier) bool {
	return identifierImplVar.name == other.getName() /*&& identifierImplVar.alias == other.GetAlias()*/
}

func (identifierImplVar *identifierImpl) String() string {
	// return "[" + identifierImplVar.name + "," + identifierImplVar.alias + "]"
	return "[" + identifierImplVar.name + "]"
}
