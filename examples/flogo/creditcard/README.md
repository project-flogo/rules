## Rules example for processing Creditcard application

### Description

This example demonstrates rule based processing of credit card application considering age and credit score. We will be using two tuples. Tuple n1 with user details and tuple n2 with credit score details.

Rule 1 :

Conditions : 
* Address should not be empty
* User age >= 18 years

Action :
* Application submitted

Rule 2 : 

Conditions :
* Match with applicant id and creditscore id
* Match with applicant name and creditscore name 
* User creditscore >= 750 
* Eligible creditlimit condition check

Action :
* Credit card application approved

If condition satisfies then respective actions gets triggered.
This example is implemented in api and json methods.

## API method

* Execute
``` 
$ cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard/api
$ go run main.go
```
## JSON method
* Create flogo app
```
$ cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard/json
$ flogo create -f flogo.json
```
* copy functions.go to cardapp/src folder
```
cp functions.go cardapp/src
```
* Run flogo build
* Start flogo app
* Input user details using below curl request

```
$ curl -X PUT http://localhost:7777/applicant -H 'Content-Type: application/json' -d '{"name":"Sam4","age":"26","salary":"50100","address":"SFO","id":"4"}'
```
* Input credit details of the user

```
$ curl -X PUT http://localhost:7777/cibildata -H 'Content-Type: application/json' -d '{"id":"4","name":"Sam4","creditScore":"850"}'
```

* Application status will be printed on the console
 