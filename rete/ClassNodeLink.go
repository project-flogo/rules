package rete

//classNodeLink links the classNode to the rest of the Rule's network
type classNodeLink interface {
	nodeLink
	getIdentifier() identifier
	getClassNode() classNode
}

type classNodeLinkImpl struct {
	nodeLinkImpl
	rule          Rule
	identifierVar identifier
	classNodeVar  classNode
}

func newClassNodeLink(classNodeVar classNode, child node, rule Rule, identifierVar identifier) classNodeLink {
	classNodeLinkImplVar := classNodeLinkImpl{}
	classNodeLinkImplVar.initClassNodeLinkImpl(classNodeVar, child, rule, identifierVar)
	return &classNodeLinkImplVar
}

func (classNodeLinkImplVar *classNodeLinkImpl) initClassNodeLinkImpl(classNodeVar classNode, child node, rule Rule, identifierVar identifier) {
	initClassNodeLink(&classNodeLinkImplVar.nodeLinkImpl, child)
	classNodeLinkImplVar.classNodeVar = classNodeVar
	classNodeLinkImplVar.rule = rule
	classNodeLinkImplVar.identifierVar = identifierVar
}

func (classNodeLinkImplVar *classNodeLinkImpl) getIdentifier() identifier {
	return classNodeLinkImplVar.identifierVar
}

func (classNodeLinkImplVar *classNodeLinkImpl) getClassNode() classNode {
	return classNodeLinkImplVar.classNodeVar
}

func (classNodeLinkImplVar *classNodeLinkImpl) String() string {
	str := classNodeLinkImplVar.nodeLinkImpl.String()
	return str
	//TODO: tableids, loadstop, mask, etc
}
