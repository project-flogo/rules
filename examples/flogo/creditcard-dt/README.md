# Decision Table Usage

This example demonstrates how to use decision table activity with credit card application example.

## Setup and build
Once you have the `flogo.json` file, you are ready to build your Flogo App

### Pre-requisites
* Go 1.11
* Download and build the Flogo CLI 'flogo' and add it to your system PATH

### Steps

```sh
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard-dt
flogo create -f flogo.json
cd creditcard-dt
flogo build
cd bin
./creditcard-dt
```
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

### Writing Decision Table in JSON

Sample usage can be as below.
```json
 {
            "name": "ApplicantSimple",
            "description": "Simple Applicants approval dt",
            "type": "activity",
            "ref": "github.com/project-flogo/rules/activity/dtable",
            "settings": {
              "make": [
                {
                  "condition": [
                    {"tuple": "applicant","field": "name","expr": "== 'JohnDoe'"},
                    {"tuple": "applicant","field": "age","expr": ">= 20"},
                    {"tuple": "applicant","field": "age","expr": "<= 30"}
                  ],
                  "action": [
                   
                    { "tuple": "applicant","field": "creditLimit","value": 2500.0},
                    { "tuple": "applicant","field": "eligible","value": true},
                    { "tuple": "applicant","field": "status","value": "VISA-Granted"}
                  ]
                },
                {
                  "condition": [
                    {"tuple": "applicant","field": "name","expr": "== 'SandraW'"},
                    {"tuple": "applicant","field": "age","expr": ">= 20"},
                    {"tuple": "applicant","field": "age","expr": "<= 30"}
                  ],
                  "action": [
                    
                    { "tuple": "applicant","field": "creditLimit","value": 0.0},
                    { "tuple": "applicant","field": "eligible","value": false},
                    { "tuple": "applicant","field": "status","value": "Loan-Rejected"}
                  ]
                },
                {
                  "condition": [
                    {"tuple": "applicant","field": "name","expr": "== 'PrakashY'"},
                    {"tuple": "applicant","field": "age","expr": ">= 20"},
                    {"tuple": "applicant","field": "age","expr": "<= 30"}
                  ],
                  "action": [
                    
                    { "tuple": "applicant","field": "creditLimit","value": 7500.0},
                    { "tuple": "applicant","field": "eligible","value": true},
                    { "tuple": "applicant","field": "status","value": "Pending"}
                  ]
                },
                {
                  "condition": [
                    {"tuple": "applicant","field": "name","expr": "== 'JaneDoe'"},
                    {"tuple": "applicant","field": "age","expr": ">30"}
                  ],
                  "action": [
                    
                    { "tuple": "applicant","field": "creditLimit","value": 25000.0},
                    { "tuple": "applicant","field": "eligible","value": false},
                    { "tuple": "applicant","field": "status","value": "Platinum-Status"}
                  ]
                }
              ]
            }
          }
```
Decision table will have condition and action included into decition table activity.