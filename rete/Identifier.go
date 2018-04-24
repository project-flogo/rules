package rete

//Identifier An internal representation of a 'DataSource'
type Identifier interface {
	getName() string
	equals(Identifier) bool
	String() string
}

type identifierImpl struct {
	name string
	// alias string
}

func newIdentifier(name string /* alias string */) Identifier {
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

func (identifierImplVar *identifierImpl) equals(other Identifier) bool {
	return identifierImplVar.name == other.getName() /*&& identifierImplVar.alias == other.GetAlias()*/
}

func (identifierImplVar *identifierImpl) String() string {
	// return "[" + identifierImplVar.name + "," + identifierImplVar.alias + "]"
	return "[" + identifierImplVar.name + "]"
}
