# Decision Table Usage

This example demonstrates how to use decision table activity with credit card application example.

## Setup and build
Once you have the `flogo.json` file, you are ready to build your Flogo App

### Pre-requisites
* Go 1.11
* Download and build the Flogo CLI 'flogo' and add it to your system PATH

### Steps

Note: Store implementation can be configured via given `rsconfig.json` file. Start redis-server and use `export STORECONFIG=<path to rsconfig.json>` before running binary.<br>

```sh
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard-dt
flogo create -f flogo.json
cd creditcard-dt
flogo build
cd bin
cp ../../creditcard-dt-file.xlsx .
./creditcard-dt
```
### Note - To execute following commands on windows platform the  backward slashes need to be omitted . 
### Testing

#### #1 Invoke applicant decision table

Store aplicants information.
```sh
curl localhost:7777/test/applicant?name=JohnDoe\&gender=Male\&age=20\&address=BoltonUK\&hasDL=false\&ssn=1231231234\&income=45000\&maritalStatus=single\&creditScore=500
curl localhost:7777/test/applicant?name=JaneDoe\&gender=Female\&age=38\&address=BoltonUK\&hasDL=false\&ssn=2424354532\&income=32000\&maritalStatus=single\&creditScore=650
curl localhost:7777/test/applicant?name=PrakashY\&gender=Male\&age=30\&address=RedwoodShore\&hasDL=true\&ssn=2345342132\&income=150000\&maritalStatus=married\&creditScore=750
curl localhost:7777/test/applicant?name=SandraW\&gender=Female\&age=26\&address=RedwoodShore\&hasDL=true\&ssn=3213214321\&income=50000\&maritalStatus=single\&creditScore=625
```

Send a process application event.
```sh
curl localhost:7777/test/processapplication?start=true\&ssn=1231231234
curl localhost:7777/test/processapplication?start=true\&ssn=2345342132
curl localhost:7777/test/processapplication?start=true\&ssn=3213214321
curl localhost:7777/test/processapplication?start=true\&ssn=2424354532
```
You should see following output:
```
2019-09-24T12:54:08.674+0530    INFO    [flogo.rules] -  Applicant: JohnDoe -- CreditLimit: 2500 -- status: VISA-Granted
2019-09-24T12:54:08.683+0530    INFO    [flogo.rules] -  Applicant: PrakashY -- CreditLimit: 7500 -- status: Pending
2019-09-24T12:54:08.696+0530    INFO    [flogo.rules] -  Applicant: SandraW -- CreditLimit: 0 -- status: Loan-Rejected
2019-09-24T12:54:09.884+0530    INFO    [flogo.rules] -  Applicant: JaneDoe -- CreditLimit: 25000 -- status: Platinum-Status
```
