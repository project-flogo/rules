package rete

//classNodeLink links the classNode to the rest of the Rule's network
type classNodeLink interface {
	nodeLink
	getIdentifier() Identifier
	getClassNode() classNode
}

type classNodeLinkImpl struct {
	nodeLinkImpl
	rule         Rule
	identifier   Identifier
	classNodeVar classNode
}

func newClassNodeLink(classNodeVar classNode, child node, rule Rule, identifier Identifier) classNodeLink {
	classNodeLinkImplVar := classNodeLinkImpl{}
	classNodeLinkImplVar.initClassNodeLinkImpl(classNodeVar, child, rule, identifier)
	return &classNodeLinkImplVar
}

func (classNodeLinkImplVar *classNodeLinkImpl) initClassNodeLinkImpl(classNodeVar classNode, child node, rule Rule, identifier Identifier) {
	initClassNodeLink(&classNodeLinkImplVar.nodeLinkImpl, child)
	classNodeLinkImplVar.classNodeVar = classNodeVar
	classNodeLinkImplVar.rule = rule
	classNodeLinkImplVar.identifier = identifier
}

func (classNodeLinkImplVar *classNodeLinkImpl) getIdentifier() Identifier {
	return classNodeLinkImplVar.identifier
}

func (classNodeLinkImplVar *classNodeLinkImpl) getClassNode() classNode {
	return classNodeLinkImplVar.classNodeVar
}

func (classNodeLinkImplVar *classNodeLinkImpl) String() string {
	str := classNodeLinkImplVar.nodeLinkImpl.String()
	return str
	//TODO: tableids, loadstop, mask, etc
}
