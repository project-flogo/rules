# BE*Go*

BE*Go* is a lightweight rules engine written in Go.

## Definitions
A `StreamSource` represents a streaming data source a.k.a Channel

A `StreamTuple` is an instance tuple of a certain `StreamSource`

A `Rule` constitutes of multiple Conditions and the rule trigers when all its conditions pass

A `Condition` is an expression comprising of data tuples from one or more StreamSources. When the expression evaluates to true, the condition passes. Thus a Condition is characterized by the number and types of StreamSource involved in its evaluation. In order to optimize a Rule's evaluation, the Rule network needs to know of the number and type of stream sources in each of its Conditions. Thus the Condition can either be defined upfront with the number of and type of StreamSources or this can be inferred (for a start, we choose the former)

The rule's `Action` is a function that is invoked once each time that a matching combination of tuples is found that result in a true evaluation of all its conditions. Thus the `Action` function takes as an argument, the matching tuples.

Each Rule creates its own evaluation plan or a network. When there are multiple such rules, all of them collectively form the Rule Network

When streaming data arrives, a StreamTuple is formed, and is `Assert`-ed into this network to see which rules would fire.
Similarly, removing a StreamTuple from the network is called `Retract`-ing it from the network. You may retract when for example, the TTL of a StreamTuple expires based on a certain expiration policy

## Server API
With this backgroud, let us see how it translates to API/code. *This is a draft Server-side API*


	//Create Rule with a name
	rule := ruleapi.NewRule("My first rule")

    //Add a condition named c1 o the rule, saying that this condition needs data from three stream sources called n1, n2 and n3 and the go function to evaluate this condition is myCondition
    //(see the myCondition signature below)
	rule.AddCondition("c1", []model.StreamSource{"n1", "n2", "n3"}, myCondition)

    //Similarly, add another named c2 to the rule
    rule.AddCondition("c2", []model.StreamSource{"n1", "n2"}, myCondition)

    //Add an Action callback function myActionFn to the Rule
	rule.SetActionFn(myActionFn)

	//Create a RuleSession. All interactions happen via this session
	ruleSession := ruleapi.NewRuleSession()

    //Add the rule to the session, you can add multiple rules such as above to the session
	ruleSession.AddRule(rule)

	//Simulate/create a few StreamTuples
	streamTuple1 := model.NewStreamTuple("n1") //simulate a new tuple of type n1
    streamTuple2 := model.NewStreamTuple("n2") //simulate a new tuple of type n2
	streamTuple3 := model.NewStreamTuple("n3") //simulate a new tuple of type n3

    //Assert them into the session
	ruleSession.Assert(streamTuple1)
	ruleSession.Assert(streamTuple2)
    ruleSession.Assert(streamTuple3)

    //Retract them
    ruleSession.Retract (streamTuple1)
    ruleSession.Retract (streamTuple2)
    ruleSession.Retract (streamTuple3)

    //You may remove the rule
    ruleSession.DeleteRule (rule.getName())


    func myCondition(ruleName string, condName string, tuples map[model.StreamSource]model.StreamTuple)     bool {
	    fmt.Printf("Condition [%s] of Rule [%s] has [%d] tuples\n", condName, ruleName, len(tuples))
	    return true
    }

    func myActionFn(ruleName string, tuples map[model.StreamSource]model.StreamTuple) {
	    fmt.Printf("My Rule [%s] fired\n", ruleName)
    }

## Try it out
* Check out the repo say at `/home/yourname/go` and set environment variable GOPATH to it
* Goto github.com/TIBCOSoftware/bego such that your path looks like this
* /home/yourname/go/src/github.com/TIBCOSoftware/bego
* Go to that the folder above and 
* `go install ./...`
* This will create an example executable at `/home/yourname/go/bin/rulesapp`
* `cd /home/yourname/go/bin`
* `./rulesapp`


## Todo
* Lots of stuff, including multi-threading support
* Review and refine the API
* Client API (if required)
* Transports/Channels (as required)
* Forward-chaining/Conflict Resolution/Post-RTC etc
* Timers and expiry

