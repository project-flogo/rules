## RuleAction
RuleAction is a flogo action which allows events to be injected into the rules engine via a flogo trigger.

### Configuration
RuleAction configuration contains 2 parts, settings and resources.

#### settings

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| id | string | id is referenced by an element in another section of flogo configuration such as trigger handler action's id |
| rulesessionURI | uri | Uri that starts with 'res://rulesession:'. It's referenced in the resources section  |
| tds | array | Tuple definitions |


#### tds
| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| name | string | Tuple type name |
| properties | array | Properties of the tuple |


#### properties
| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| name | string | Property name |
| type | string | Data type of the property |
| pk-index | int | Tuple key order. If -1, not used as part of tuple key |

#### resources

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| id | string | id is referenced by an element in another section of flogo configuration such as action settings's rulesessionURI |
| data | object | metadata and rule defintions |

#### metadata
It contains an array of input element comprised of 2 parameters, values and tupletype.

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| values | string | Tuple values |
| tupletype | string | Tupple type that the values parameter adheres to |

#### rules
| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| name | string | Name of the rule |
| conditions | array | Conditions that the rule evaluates given input |
| actionFunction | string | Rule action function to be fired when conditions are true. The function must exist in functions.go |

#### conditions

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| name | string | Name of the condition |
| identifiers | array | Tuple types the condition evaluates upon |
| evaluator | string | Function that envaluates the condition. The function must exist in functions.go |


### Usage

Upon a configuration of flogo.json containing ruleaction configuration and creation of functions.go which contains the evaluator and actionFunction implementations, place flogo.json in a folder of your choice to run 
```
flogo create -f flogo.json myrules
cp functions.go myrules/src
cd myrules
flogo build
```
To run the flogo rules binary,
```
bin/myrules
```


### Examples
For examples, see [rules flogo examples](https://github.com/project-flogo/rules/tree/master/examples/flogo).