# Decision Table Usage

This example demonstrates how to use decision table activity with student analysis example.

## Setup and build
Once you have the `flogo.json` file, you are ready to build your Flogo App

### Pre-requisites
* Go 1.11
* Download and build the Flogo CLI 'flogo' and add it to your system PATH

### Steps

```sh
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/dtable
flogo create -f flogo.json
cd dtable
flogo build
cd bin
./dtable
```
### Testing

#### #1 Invoke student analysis decision table

Store student information.
```sh
curl localhost:7777/test/student?grade=GRADE-C\&name=s1\&class=X-A\&careRequired=false
curl localhost:7777/test/student?grade=GRADE-B\&name=s2\&class=X-A\&careRequired=false
```

Send a curl student analysis event.
```sh
curl localhost:7777/test/studentanalysis?name=s1
curl localhost:7777/test/studentanalysis?name=s2

```
You should see following output:
```
2019-09-05T18:35:12.142+0530    INFO    [flogo.rules] - Student: s1 -- Comments: additional study hours required
2019-09-05T18:35:12.142+0530    INFO    [flogo.rules] - Student: s2 -- Comments: little care can be taken to achieve grade-a
```

### Writing Decision Table in JSON

Sample usage can be as below.
```json
    {
            "name": "AnalyseStudent",
            "description": "Analysing student data",
            "type": "activity",
            "ref": "github.com/project-flogo/rules/activity/dtable",
            "settings": {
              "make": [
                {
                  "condition": [
                    {"tuple": "student","field": "grade","expr": "== 'GRADE-C'"},
                    {"tuple": "student","field": "class","expr": "== 'X-A'"}
                  ],
                  "action": [
                    { "tuple": "student","field": "careRequired","value": true},
                    { "tuple": "student","field": "comments","value": "additional study hours required"}
                  ]
                },
                {
                  "condition": [
                    {"tuple": "student","field": "grade","expr": "== 'GRADE-A'"},
                    {"tuple": "student","field": "class","expr": "== 'X-A'"}
                  ],
                  "action": [
                    {"tuple": "student","field": "careRequired","value": false}
                  ]
                },
                {
                  "condition": [
                    {"tuple": "student","field": "grade","expr": "== 'GRADE-B'"},
                    {"tuple": "student","field": "class","expr": "== 'X-A'"}
                  ],
                  "action": [
                    {"tuple": "student","field": "careRequired","value": true},
                    {"tuple": "student","field": "comments","value": "little care can be taken to achieve grade-a "}
                  ]
                }
              ]
            }
        }
```
Decision table will have condition and action included into decition table activity.