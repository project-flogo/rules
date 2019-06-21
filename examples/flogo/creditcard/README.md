## Flogo Rules based Creditcard application


This example demonstrates rule based processing of credit card application. In this example three tuples are used, tuples description is given below.


* `UserAccount` tuple is always stored in network, while the other tuples `NewAccount` and `UpdateCibil` are removed after usage as ttl is given as 30secs and 0. 
* During startup Name Tom is asserted into network.

## Usage

Get the repo and in this example main.go, functions.go both are available. We can directly build and run the app or create flogo rule app and run it.

#### Conditions 

```
cUserData : Check tuple not empty
cNewUser : Check for new user input data - checks if age <17 and >=45, empty address and salary less than 10k
cNewUserId : Check for id match from 'UserAccount' and 'NewAccount' tuples
cNewUserAge : Check for age >=18 and <= 44 
cNewUserCibil : Check for cibil >= 750 && < 820 
cNewUserLowCibil : Check for cibil <750
cNewUserHighCibil : Check for cibil >= 820 &&  <= 900
```
#### Actions 
```
aUserData : Executes when Proper Input details provided
aNewUser : When required values for addrees, age, salay not matching with pre-requisites and retracts NewAccount tuple
aNewUserApprove : Provides credit card application status approved with prescribed credit limit
aNewUserApprove1 : Provides credit card application status approved with prescribed credit limitt tuple
aNewUserReject : Rejects when lower cibil score provided and retracts NewAccount
```
### Direct build and run
```
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard
go build
./creditcard
```
### Create app using flogo cli
```
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard
flogo create -f flogo.json creditcard
cp functions.go creditcard/src
cd creditcard
flogo build
./bin/creditcard
```

* Input user details - Stored in network

```
$ curl -X PUT http://localhost:7777/users -H 'Content-Type: application/json' -d '{"Id":"12312","Name":"Tom1","Age":19,"Addres":"TEST","Gender":"male","maritalStatus":"single","appStatus":""}'
```
* Input new user details - Stored in network for 30 secs

```
$ curl -X PUT http://localhost:7777/newaccount -H 'Content-Type: application/json' -d '{"Name":"Test","Age":"26","Income":"60100","Address":"TEt","Id":"12312","Gener":"male","maritalStatus":"single"}'
```
* Update credit score details of the user

```
$ curl -X PUT http://localhost:7777/credit -H 'Content-Type: application/json' -d '{"Id":12312,"creditScore":680}'
```

* Application status will be printed on the console
 