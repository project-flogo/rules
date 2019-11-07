# Decision Table
This is a `flogo activity` based `Decision Table` implementation used inside a rules application. It can be invoked as rule action service. `Decision Table` provide a tabular way to build complex business rules. Each column can be created based on predefined properties. These properties are defined in the `tuple descriptor` used inside a rule application. Each row can be thought of as one rule in a table made up of many rules. The individual rules are often straightforward, as in the following example.

Rule conditions:
```
person.age > 30
person.gender == "male"
```
Rule actions:
```
application.status = "ACCEPTED"
application.credit = 4000
```
A decision table can consist of hundreds, even thousands of rules each of which is executed only when its specific conditions are satisfied.

## Usage

The available activity `settings` are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| dTableFile | string |  Decision table file path (xlsx & csv extensions are supported) |

A sample `decision table` definition is:
```json
    {
        "name": "ApplicantSimple",
        "description": "Simple Applicants approval dt",
        "type": "activity",
        "ref": "github.com/project-flogo/rules/activity/dtable",
        "settings": {
            "dTableFile":"creditcard-dt-file.xlsx"
        }
    }

```

An example rule action service that invokes the above `ApplicantSimple` decision table is:
```json
    "actionService": {
        "service": "ApplicantSimple",
        "input": {
            "message": "test ApplicantSimple"
        }
    }
```