#!/bin/bash

function get_test_cases {
    local my_list=( testcase1 testcase2 testcase3 )
    echo "${my_list[@]}"
}

# Test cases performs credit card application status as approved if Creditscore > 750
function testcase1 {
pushd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard/json
flogo create -f flogo.json
cp functions.go cardapp/src
cd cardapp
flogo build
./bin/cardapp > /tmp/testcase1.log 2>&1 &
pId=$!

response=$(curl -X PUT http://localhost:7777/applicant -H 'Content-Type: application/json' -d '{"name":"Sam4","age":"26","salary":"50100","address":"SFO","id":"4"}' --write-out '%{http_code}' --silent --output /dev/null)
response1=$(curl -X PUT http://localhost:7777/CreditScoredata -H 'Content-Type: application/json' -d '{"id":"4","name":"Sam4","creditScore":"850"}' --write-out '%{http_code}' --silent --output /dev/null)
kill -9 $pId

if [ $response -eq 200 ] && [ $response1 -eq 200 ] && [[ "echo $(cat /tmp/testcase1.log)" =~ "Rule fired" ]]
    then 
        echo "PASS"
    else
        echo "FAIL"
fi
cd ..
rm -rf cardapp
popd
}

# Test cases performs credit card application status rejected if Creditscore < 750
function testcase2 {
pushd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard/json
flogo create -f flogo.json
cp functions.go cardapp/src
cd cardapp
flogo build
./bin/cardapp > /tmp/testcase2.log 2>&1 &
pId=$!

response=$(curl -X PUT http://localhost:7777/applicant -H 'Content-Type: application/json' -d '{"name":"Sam4","age":"26","salary":"50100","address":"SFO","id":"4"}' --write-out '%{http_code}' --silent --output /dev/null)
response1=$(curl -X PUT http://localhost:7777/CreditScoredata -H 'Content-Type: application/json' -d '{"id":"4","name":"Sam4","creditScore":"650"}' --write-out '%{http_code}' --silent --output /dev/null)

kill -9 $pId
if [ $response -eq 200  ] && [ $response1 -eq 200  ] && [[ "echo $(cat /tmp/testcase2.log)" =~ "lower CreditScore score" ]]
    then 
        echo "PASS"
    else
        echo "FAIL"
fi
cd ..
rm -rf cardapp
popd
}

# Test cases performs credit card application status as approved if Creditscore > 750
function testcase3 {
pushd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard/api
go run main.go > /tmp/testcase3.log 2>&1

if [[ "echo $(cat /tmp/testcase3.log)" =~ "Rule fired" ]]
    then 
        echo "PASS"
    else
        echo "FAIL"
fi
popd
}