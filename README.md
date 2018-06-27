# BE*Go*

BE*Go* is a lightweight rules engine written in Go.

## Definitions
A `StreamSource` represents a streaming data source a.k.a Channel

A `StreamTuple` is an instance of a certain type of `StreamSource`

A `Rule` constitutes of multiple Conditions and the rule triggers when all its conditions pass

A `Condition` is an expression comprising of data tuples from one or more StreamSources. When the expression evaluates to true, the condition passes. Thus a Condition is characterized by the number and types of StreamSource involved in its evaluation. In order to optimize a Rule's evaluation, the Rule network needs to know of the number and type of stream sources in each of its Conditions. Thus the Condition can either be defined upfront with the number of and type of StreamSources or this can be inferred (for a start, we choose the former)

The rule's `Action` is a function that is invoked once each time that a matching combination of tuples is found that result in a true evaluation of all its conditions. Thus the `Action` function takes as an argument, the matching tuples.

Each Rule creates its own evaluation plan or a network. When there are multiple such rules, all of them collectively form the Rule Network

When streaming data arrives, a StreamTuple is formed, and is `Assert`-ed into this network to see which rules would fire.
Similarly, removing a StreamTuple from the network is called `Retract`-ing it from the network. You may retract when for example, the TTL of a StreamTuple expires based on a certain expiration policy

## Server API
With this background, let us see how it translates to API/code. *This is a draft Server-side API*


	//Create Rule, define conditions and set action callback
	rule := ruleapi.NewRule("* Ensure n1.name is Bob and n2.name matches n1.name ie Bob in this case *")
	fmt.Printf("Rule added: [%s]\n", rule.GetName())
	rule.AddCondition("c1", []model.StreamSource{"n1"}, checkForBob)          // check for name "Bob" in n1
	rule.AddCondition("c2", []model.StreamSource{"n1", "n2"}, checkSameNames) // match the "name" field in both tuples
	//in effect, fire the rule when name field in both tuples is "Bob"
	rule.SetActionFn(myActionFn)

	//Create a RuleSession and add the above Rule
	ruleSession := ruleapi.NewRuleSession()
	ruleSession.AddRule(rule)

	//Now assert a few facts and see if the Rule Action callback fires.
	fmt.Println("Asserting n1 tuple with name=Bob")
	streamTuple1 := model.NewStreamTuple("n1")
	streamTuple1.SetString("name", "Bob")
	ruleSession.Assert(streamTuple1)

	fmt.Println("Asserting n1 tuple with name=Fred")
	streamTuple2 := model.NewStreamTuple("n1")
	streamTuple2.SetString("name", "Fred")
	ruleSession.Assert(streamTuple2)

	fmt.Println("Asserting n2 tuple with name=Fred")
	streamTuple3 := model.NewStreamTuple("n2")
	streamTuple3.SetString("name", "Fred")
	ruleSession.Assert(streamTuple3)

	fmt.Println("Asserting n2 tuple with name=Bob")
	streamTuple4 := model.NewStreamTuple("n2")
	streamTuple4.SetString("name", "Bob")
	ruleSession.Assert(streamTuple4)

    //Retract them
    ruleSession.Retract (streamTuple1)
    ruleSession.Retract (streamTuple2)
    ruleSession.Retract (streamTuple3)
    ruleSession.Retract (streamTuple4)

    //You may delete the rule
    ruleSession.DeleteRule (rule.getName())

    func checkForBob(ruleName string, condName string, tuples map[model.StreamSource]model.StreamTuple) bool {
        //This conditions filters on name="Bob"
        streamTuple := tuples["n1"]
        if streamTuple == nil {
            fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
            return false
        }
        name := streamTuple.GetString("name")
        return name == "Bob"
    }
    
    func checkSameNames(ruleName string, condName string, tuples map[model.StreamSource]model.StreamTuple) bool {
        // fmt.Printf("Condition [%s] of Rule [%s] has [%d] tuples\n", condName, ruleName, len(tuples))
        streamTuple1 := tuples["n1"]
        streamTuple2 := tuples["n2"]
        if streamTuple1 == nil || streamTuple2 == nil {
            fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
            return false
        }
        name1 := streamTuple1.GetString("name")
        name2 := streamTuple2.GetString("name")
        return name1 == name2
    }
    
    func myActionFn(ruleName string, tuples map[model.StreamSource]model.StreamTuple) {
        fmt.Printf("Rule fired: [%s]\n", ruleName)
        streamTuple1 := tuples["n1"]
        streamTuple2 := tuples["n2"]
        if streamTuple1 == nil || streamTuple2 == nil {
            fmt.Println("Should not get nil tuples here in Action! This is an error")
        }
        name1 := streamTuple1.GetString("name")
        name2 := streamTuple2.GetString("name")
        fmt.Printf("n1.name = [%s], n2.name = [%s]\n", name1, name2)
    }

## Try it out
* Checkout the repo at say `/home/yourname/go` and set environment variable `$GOPATH` to it
* Goto `github.com/TIBCOSoftware/bego` such that your path looks like this
* `/home/yourname/go/src/github.com/TIBCOSoftware/bego`
* Go to the folder above and 
* `go install ./...`
* This will create the example executable at `/home/yourname/go/bin/rulesapp`
* `cd /home/yourname/go/bin`
* `./rulesapp`


## Todo
* Lots of stuff, including multi-threading support
* Review and refine the API
* Client API (if required)
* Transports/Channels (as required)
* Forward-chaining/Conflict Resolution/Post-RTC etc (TDB: if required)
* Timers and expiry
* Storage 

