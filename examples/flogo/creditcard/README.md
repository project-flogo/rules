## Flogo Rules based Creditcard application


This example demonstrates rule based processing of credit card application. In this example three tuples are used, tuples description is given below.


* `UserAccount` tuple is always stored in network, while the other tuples `NewAccount` and `UpdateCreditScore` are removed after usage as ttl is given as 0. 


## Usage

Get the repo and in this example main.go, functions.go both are available. We can directly build and run the app or create flogo rule app and run it.

#### Conditions 

```
cBadUser : Check for new user input data - check if age <18 and >=45, empty address and salary less than 10k
cNewUser : Check for new user input data - check if age >=18 and <= 44, address and salary >= 10k
cUserIdMatch : Check for id match from 'UserAccount' and 'UpdateCreditScore' tuples
cUserCreditScore : Check for CreditScore >= 750 && < 820 
cUserLowCreditScore : Check for CreditScore < 750
cUserHighCreditScore : Check for CreditScore >= 820 &&  <= 900
```
#### Actions 
```
aBadUser : Executes when age - < 18 and >=45, address empty, salary less than 10k
aNewUser : Add the newuser info to userAccount tuple
aApproveWithLowerLimit : Provides credit card application status approved with lower credit limit
aApproveWithHigherLimit : Provides credit card application status approved with higher credit limit
aUserReject : Rejects when lower Credit score provided and retracts NewAccount
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

* Input new user details

```
$ curl -X PUT http://localhost:7777/newaccount -H 'Content-Type: application/json' -d '{"Name":"Test","Age":"26","Income":"60100","Address":"TEt","Id":"12312","Gender":"male","maritalStatus":"single"}'
```
* Update credit score details of the user

```
$ curl -X PUT http://localhost:7777/credit -H 'Content-Type: application/json' -d '{"Id":12312,"creditScore":680}'
```

* Application status will be printed on the console
 