package ruleapi

type sessionCtx interface {
	setRuleSession(rs RuleSession)
	getRuleSession() RuleSession
}

type sessionKeyType struct {
}

var sessionCtxKEY = sessionKeyType{}

type sessionCtxImpl struct {
	rs RuleSession
}

func newSessionCtx() sessionCtx {
	sCtx := sessionCtxImpl{}
	return &sCtx
}

func (sctx *sessionCtxImpl) setRuleSession(rs RuleSession) {
	sctx.rs = rs
}

func (sctx *sessionCtxImpl) getRuleSession() RuleSession {
	return sctx.rs
}
