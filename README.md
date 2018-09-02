# Flogo Rules

**Flogo Rules** is a lightweight rules library written in Golang to simplify the building of Event Driven Reactive Applications. Supports Declaritive Rules, Contextual reasoning across time and space

## Installation
### Prerequisites
To get started with the Flogo Rules you'll need to have a few things
* The Go programming language version 1.8 or later should be [installed](https://golang.org/doc/install).
* The **GOPATH** environment variable on your system must be set properly

### Install
```
$ go get -u github.com/TIBCOSoftware/bego
```
_Note that the -u parameter automatically updates bego if it exists_

## Getting Started
Getting started should be fairly easy. Lets start off with some definitions around various types used.

### Definitions
A `Tuple` represents an event or a business object and provides runtime data to the rules. It is always of a certain type.

A `TupleTypeDescriptor` defines the type or structure of a `Tuple`. It defines a tuple's properties and data types, and primary keys. It also defines the time to live for the tuple

A `TupleType` is a name or an alias for a `TupleTypeDescriptor` 

A `Rule` constitutes of multiple Conditions and the rule triggers when all its conditions pass

A `Condition` is an expression involving one or more tuple types. When the expression evaluates to true, the condition passes. In order to optimize a Rule's evaluation, the Rule network needs to know of the TupleTypes and the properties of the TupleType which participate in the `Condition` evaluation. These are provided when constructing the condition and adding it to the rule.

A `Action` is a function that is invoked each time that a matching combination of tuples are found that result in a `true` evaluation of all its conditions. Those matching tuples are passed to the action function.

A `RuleSession` is a handle to interact with the rules API. You can create and register multiple rule sessions. Rule sessions are silos for the data that they hold, they are similar to namespaces. Sharing objects/state across rule sessions is not supported.
 
Each rule creates its own evaluation plan or a network. Multiple rules collectively form the rule network

Tuples can be created using `NewTuple` and then setting its properties. The tuple is then `Assert`-ed into the rule session and this triggers rule evaluations.
A tuple can be `Retract`ed from the rule session to take it out of play for rules evaluations.

### Usage
Now lets see some code in action. Below code snippet demonstrates usage of the Rules API,

	fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//First register the tuple descriptors
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	//Create a RuleSession
	rs, _ := ruleapi.GetOrCreateRuleSession("asession")

	//// check for name "Bob" in n1
	rule := ruleapi.NewRule("n1.name == Bob")
	rule.AddCondition("c1", []string{"n1"}, checkForBob, nil)
	rule.SetAction(checkForBobAction)
	rule.SetContext("This is a test of context")
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	// check for name "Bob" in n1, match the "name" field in n2,
	// in effect, fire the rule when name field in both tuples is "Bob"
	rule2 := ruleapi.NewRule("n1.name == Bob && n1.name == n2.name")
	rule2.AddCondition("c1", []string{"n1"}, checkForBob, nil)
	rule2.AddCondition("c2", []string{"n1", "n2"}, checkSameNamesCondition, nil)
	rule2.SetAction(checkSameNamesAction)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	//Now assert a "n1" tuple
	fmt.Println("Asserting n1 tuple with name=Tom")
	t1, _ := model.NewTupleWithKeyValues("n1", "Tom")
	t1.SetString(nil, "name", "Tom")
	rs.Assert(nil, t1)

	//Now assert a "n1" tuple
	fmt.Println("Asserting n1 tuple with name=Bob")
	t2, _ := model.NewTupleWithKeyValues("n1", "Bob")
	t2.SetString(nil, "name", "Bob")
	rs.Assert(nil, t2)

	//Now assert a "n2" tuple
	fmt.Println("Asserting n2 tuple with name=Bob")
	t3, _ := model.NewTupleWithKeyValues("n2", "Bob")
	t3.SetString(nil, "name", "Bob")
	rs.Assert(nil, t3)

	//Retract tuples
	rs.Retract(nil, t1)
	rs.Retract(nil, t2)
	rs.Retract(nil, t3)

	//delete the rule
	rs.DeleteRule(rule.GetName())

	//unregister the session, i.e; cleanup
	rs.Unregister()



## Try out this example
* Create a directory say `/home/rulesexample` and set `GOPATH=/home/rulesexample`
* from `/home/rulesexample` run `go get github.com/TIBCOSoftware/bego/rulesapp`
* This will create the example executable at `/home/rulesexample/bin/rulesapp`
* Run the example `/home/rulesexample/bin/rulesapp`


