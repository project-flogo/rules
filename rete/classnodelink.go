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
	cnl := classNodeLinkImpl{}
	cnl.initClassNodeLinkImpl(classNodeVar, child, rule, identifierVar)
	return &cnl
}

func (cnl *classNodeLinkImpl) initClassNodeLinkImpl(classNodeVar classNode, child node, rule Rule, identifierVar identifier) {
	initClassNodeLink(&cnl.nodeLinkImpl, child)
	cnl.classNodeVar = classNodeVar
	cnl.rule = rule
	cnl.identifierVar = identifierVar
}

func (cnl *classNodeLinkImpl) getIdentifier() identifier {
	return cnl.identifierVar
}

func (cnl *classNodeLinkImpl) getClassNode() classNode {
	return cnl.classNodeVar
}

func (cnl *classNodeLinkImpl) String() string {
	str := cnl.nodeLinkImpl.String()
	return str
	//TODO: tableids, loadstop, mask, etc
}
