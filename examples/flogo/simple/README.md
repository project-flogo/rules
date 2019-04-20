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

## Steps to configure and build a Rules based Flogo App
Below is the `flogo.json` file used in this example application. We will use this example to explain the configuration and setup of a Flogo/Rules App
```
{
  "name": "simplerules",
  "type": "flogo:app",
  "version": "0.0.1",
  "description": "Sample Flogo App",
  "appModel": "1.0.0",
  "triggers": [
    {
      "id": "receive_http_message",
      "ref": "github.com/project-flogo/contrib/trigger/rest",
      "settings": {
        "port": "7777"
      },
      "handlers": [
        {   
          "settings": {
            "method": "GET",
            "path": "/test/n1"
          },
          "actions": [
            {
              "id": "simple_rule",
              "input": {
                "tupletype": "n1",
                "values": "=$.queryParams"
              }
            }
          ]
        },
        {
          "settings": {
            "method": "GET",
            "path": "/test/n2"
          },
          "actions": [
            {
              "id": "simple_rule",
              "input": {
                "tupletype": "n2",
                "values": "=$.queryParams"
              }
            }
          ]
        }
      ]
    }
  ],
  "actions": [
    {
      "ref": "github.com/project-flogo/rules/ruleaction",
      "settings": {
        "ruleSessionURI": "res://rulesession:simple",
        "tds": [
          {
            "name": "n1",
            "properties": [
              {
                "name": "name",
                "type": "string",
                "pk-index": 0
              }
            ]
          },
          {
            "name": "n2",
            "properties": [
              {
                "name": "name",
                "type": "string",
                "pk-index": 0
              }
            ]
          }
        ]
      },
      "id": "simple_rule"
    }
  ],
  "resources": [
    {
      "id": "rulesession:simple",
      "data": {
        "metadata": {
          "input": [
            {
              "name": "values",
              "type": "string"
            },
            {
              "name": "tupletype",
              "type": "string"
            }
          ],
          "output": [
             {
               "name": "outputData",
               "type": "any"
             }
          ]
        },
        "rules": [
          {
            "name": "n1.name == Bob",
            "conditions": [
              {
                "name": "c1",
                "identifiers": [
                  "n1"
                ],
                "evaluator": "checkForBob"
              }
            ],
            "actionFunction": "checkForBobAction"
          },
          {
            "name": "n1.name == Bob \u0026\u0026 n1.name == n2.name",
            "conditions": [
              {
                "name": "c1",
                "identifiers": [
                  "n1"
                ],
                "evaluator": "checkForBob"
              },
              {
                "name": "c2",
                "identifiers": [
                  "n1",
                  "n2"
                ],
                "evaluator": "checkSameNamesCondition"
              }
            ],
            "actionFunction": "checkSameNamesAction"
          }
        ]
      }
    }
  ]
}
``` 
## Action configuration
First configure the top level `actions` section. Here, the tags `id`, `ruleSessionURI` and `tds` are user configurable.

`id` is an alias to this rule session. It is referenced in the action configuration of the triggers

`ruleSessionURI` must start with `res://rulesession:` We have given it a name `simple` This is referenced in the resources section

`tds` contains the tuple types that need to be registered with the Rules API

##Resources configuration
Under the top level `resources` section, `id` must reference a rule session configured in the actions
configuration above. Since the name of the rulessesion being referenced is `simple`, `id` takes
the value `rulesession:simple`

Under the `rules` section, declare all your rules.
`name` gives your rule a name

The `conditions` are your rule's conditions. Each condition has a `name`, a list of types being used in the condition
under `identifiers`. These identifiers should be one of `tds/names` defined in the `actions` section. The value of `evaluator` should
be an unique string. This string will bind to a Go function at runtime (explained later)
Note how you can have multiple conditions (See the second condition in the example above)

The `actionFunction` is your rule's action. It means that when the `evaluator` condition is met, invoke this function
Again, this is a unique string whose value binds to a Go function at runtime (explained later)

## Configure the trigger handler
Flogo users are perhaps already familiar with the trigger configurations. 
In this example, we configured two handlers. In the first, we have `tupletype` as `n1` and `path` as`/test/n1`
In the second, we have we have `tupletype` as `n2` and `path` as`/test/n2`
What this means is that when data arrives on URI `test/n1` we map its data to tuple type `n1` and 
when it arrives on `/test/n2` we map its data to tuple type `n2`. Note that the `tupletype` should be one of the tuple type names defined in the `tds` section

##Mapping data from the handler to tuples
To do that, we simply configure the `actions/input/values` to `$.queryParams`

##How it all comes together
When data arrives on a trigger/handler, a new `Tuple` of `handlers/actions/input/tupletype` is created
The tuple values are initialized with the HTTP query parameters and the tuple is asserted to the rules session

##Binding Go functions for actions and conditions to the string tokens defined in the descriptor
To complete the application, you need to provide the bindings for the rules' conditions and actions in the application's `main` package
In this example, look at `functions.go`
You have to provide Go functions that adhere to the condition and action function API
```//ConditionEvaluator is a function pointer for handling condition evaluations on the server side
   //i.e, part of the server side API
   type ConditionEvaluator func(string, string, map[TupleType]Tuple, RuleContext) bool
   
   //ActionFunction is a function pointer for handling action callbacks on the server side
   //i.e part of the server side API
   type ActionFunction func(context.Context, RuleSession, string, map[TupleType]Tuple, RuleContext)
``` 
You also have to code an `init()` function in the `main` package and bind your Go functions for conditions and actions the 
corresponding string tokens defined in the `resources/rules/conditions/evaluator` and `resources/rules/actionFunction` in your descriptor like so

```
//add this sample file to your flogo project
func init() {
	config.RegisterActionFunction("checkForBobAction", checkForBobAction)
	config.RegisterActionFunction("checkSameNamesAction", checkSameNamesAction)

	config.RegisterConditionEvaluator("checkForBob", checkForBob)
	config.RegisterConditionEvaluator("checkSameNamesCondition", checkSameNamesCondition)
	config.RegisterStartupRSFunction("simple", StartupRSFunction)
}
```

This completes the configuration part

## Setup and build
Once you have the `flogo.json` file and a `functions.go` file, you are ready to build your Flogo App

###Pre-requisites
* Go 1.11
* Download and build the Flogo CLI 'flogo' and add it to your system PATH

### Steps
Place `flogo.json` in a folder of your choice, then run
`flogo create -f flogo.json`
this will pull the required dependencies and create a folder with the name as defined in the top level `name` in `flogo.json`

Now `cd` into that folder, say `simplerules`

Run the following command:
`flogo install github/project-flogo/core`

Then in the `main.go` that gets created in `src`, add these lines to the import section
```
_ "github.com/project-flogo/contrib/trigger/rest"
_ "github.com/project-flogo/rules/ruleaction"
```
Add `functions.go` to the `src` folder, next to `main.go`

From the `simplerules` folder, run `flogo build`
If everything goes well, you should have an executable (in this example) at `simplerules/bin/simplerules` 

##Test your app
First, inspect the `flogo.json` and the `functions.go` to understand the rules/conditions/actions configurations
For our test app `simplerules`, from the command line,
`simplerules/bin/simplerules`
Then from another command line, send a curl request
`curl localhost:7777/test/n1?name=Bob`
You should see this o/p on the console
```
Rule fired: [n1.name == Bob]
Context is [This is a test of context]
```
and then
`curl localhost:7777/test/n2?name=Bob`
and you should see this on the console
```
Rule fired: [n1.name == Bob && n1.name == n2.name]
n1.name = [Bob], n2.name = [Bob]
```
