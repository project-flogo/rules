package rete

import "github.com/TIBCOSoftware/bego/common/model"

//classNodeLink links the classNode to the rest of the model.Rule's network
type classNodeLink interface {
	nodeLink
	GetIdentifier() model.TupleTypeAlias
	getClassNode() classNode
	getRule() model.Rule
}

type classNodeLinkImpl struct {
	nodeLinkImpl
	rule          model.Rule
	identifierVar model.TupleTypeAlias
	classNodeVar  classNode
}

func newClassNodeLink(nw Network, classNodeVar classNode, child node, rule model.Rule, identifierVar model.TupleTypeAlias) classNodeLink {
	cnl := classNodeLinkImpl{}
	cnl.initClassNodeLinkImpl(nw, classNodeVar, child, rule, identifierVar)
	return &cnl
}

func (cnl *classNodeLinkImpl) initClassNodeLinkImpl(nw Network, classNodeVar classNode, child node, rule model.Rule, identifierVar model.TupleTypeAlias) {
	initClassNodeLink(nw, &cnl.nodeLinkImpl, child)
	cnl.classNodeVar = classNodeVar
	cnl.rule = rule
	cnl.identifierVar = identifierVar
}

func (cnl *classNodeLinkImpl) GetIdentifier() model.TupleTypeAlias {
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