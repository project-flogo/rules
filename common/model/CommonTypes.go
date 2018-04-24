package model

//ConditionEvaluator is a function pointer for handling condition evaluations on the server side
//i.e, part of the server side API
type ConditionEvaluator func(string, string, map[StreamSource]StreamTuple) bool

//ActionFunction is a function pointer for handling action callbacks on the server side
//i.e part of the server side API
type ActionFunction func(string, map[StreamSource]StreamTuple)
