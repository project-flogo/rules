package rete

import "github.com/project-flogo/rules/common/model"

//classNodeLink links the classNode to the rest of the model.Rule's network
type classNodeLink interface {
	nodeLink
	GetIdentifier() model.TupleType
	getClassNode() classNode
	getRule() model.Rule
}

type classNodeLinkImpl struct {
	nodeLinkImpl
	rule          model.Rule
	identifierVar model.TupleType
	classNodeVar  classNode
}

func newClassNodeLink(nw Network, classNodeVar classNode, child node, rule model.Rule, identifierVar model.TupleType) classNodeLink {
	cnl := classNodeLinkImpl{}
	cnl.initClassNodeLinkImpl(nw, classNodeVar, child, rule, identifierVar)
	return &cnl
}

func (cnl *classNodeLinkImpl) initClassNodeLinkImpl(nw Network, classNodeVar classNode, child node, rule model.Rule, identifierVar model.TupleType) {
	cnl.initClassNodeLink(nw, child)
	cnl.classNodeVar = classNodeVar
	cnl.rule = rule
	cnl.identifierVar = identifierVar
}

func (cnl *classNodeLinkImpl) GetIdentifier() model.TupleType {
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

func (cnl *classNodeLinkImpl) getRule() model.Rule {
	return cnl.rule
}
