<p align="center">
  <img src ="https://raw.githubusercontent.com/TIBCOSoftware/flogo/master/images/flogo-ecosystem_Rules.png" />
</p>

<p align="center" >
  <b>Rules is a lightweight library written in Golang to simplify the building of contextually aware, declaritive rules.</b>
</p>

<p align="center">
  <img src="https://travis-ci.org/TIBCOSoftware/flogo.svg"/>
  <img src="https://img.shields.io/badge/dependencies-up%20to%20date-green.svg"/>
  <img src="https://img.shields.io/badge/license-BSD%20style-blue.svg"/>
  <a href="https://gitter.im/project-flogo/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link"><img src="https://badges.gitter.im/Join%20Chat.svg"/></a>
</p>

## Installation

### Prerequisites
To get started with the Flogo Rules you'll need to have a few things
* The Go programming language version 1.8 or later should be [installed](https://golang.org/doc/install).
* The **GOPATH** environment variable on your system must be set properly

### Install
```
$ go get -u github.com/project-flogo/rules/...
```
_Note that the -u parameter automatically updates rules if it exists_

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

First we start off with loading the `TupleDescriptor`. It accepts a JSON string defining all the tuple descriptors.

	fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//First register the tuple descriptors
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

Next create a `RuleSession` and add all the `Rule`s with their `Condition`s and `Actions`s.

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
	
	//Finally, start the rule session before asserting tuples
	//Your startup function, if registered will be invoked here
	rs.Start(nil)

Here we create and assert the actual `Tuple's` which will be evaluated against the `Rule's` `Condition's` defined above.

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

Finally, once all `Rule` `Condition's` are evaluated and `Action's` are executed, we can `Retract` all the `Tuple's` from the `RuleSession` and unregister the RuleSession.

	//Retract tuples
	rs.Retract(nil, t1)
	rs.Retract(nil, t2)
	rs.Retract(nil, t3)

	//delete the rule
	rs.DeleteRule(rule.GetName())

	//unregister the session, i.e; cleanup
	rs.Unregister()

### Try out this example

```
$ go get github.com/project-flogo/rules/examples/rulesapp
```
Either manually run from source
```
$ cd $GOPATH/src/github.com/project-flogo/rules/examples/rulesapp
$ go run main.go
```
or install and run

```
$ cd $GOPATH/src/github.com/project-flogo/rules/examples/rulesapp
$ go install
$ ./$GOPATH/bin/rulesapp

```
## Running Rules in a Flogo App
To use the Rules action in your Flogo App, refer to `examples/flogo/simple/README`

## Connect with us

If you have any questions, feel free to post an issue and tag it as a question, email flogo-oss@tibco.com or chat with the team and community:

* The [project-flogo/Lobby](https://gitter.im/project-flogo/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link) Gitter channel should be used for general discussions, start here for all things Flogo/Flogo Rules,etc!
* The [project-flogo/developers](https://gitter.im/project-flogo/developers?utm_source=share-link&utm_medium=link&utm_campaign=share-link) Gitter channel should be used for developer/contributor focused conversations.

## License 
Flogo Rules source code in [this](https://github.com/project-flogo/rules) repository is under a BSD-style license, refer to [LICENSE](https://github.com/project-flogo/rules/blob/master/LICENSE) 
