package rete

//common interface for class, filter and join nodes, required so that we can call Stringer.String on them all
type abstractNode interface {
	String() string
}
