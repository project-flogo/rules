#!/bin/bash

function get_test_cases {
    local my_list=( testcase1 testcase2 testcase3 )
    echo "${my_list[@]}"
}

# Test cases performs credit card application status as approved if Creditscore > 750
function testcase1 {
pushd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard
go build 
./creditcard > /tmp/testcase1.log 2>&1 &
pId=$!

response=$(curl -X PUT http://localhost:7777/newaccount -H 'Content-Type: application/json' -d '{"Name":"Sam4","Age":"26","Income":"50100","Address":"SFO","Id":"4"}' --write-out '%{http_code}' --silent --output /dev/null)
response1=$(curl -X PUT http://localhost:7777/credit -H 'Content-Type: application/json' -d '{"Id":"4","creditScore":"850"}' --write-out '%{http_code}' --silent --output /dev/null)
kill -9 $pId

if [ $response -eq 200 ] && [ $response1 -eq 200 ] && [[ "echo $(cat /tmp/testcase1.log)" =~ "Rule fired" ]]
    then 
        echo "PASS"
    else
        echo "FAIL"
fi
cd ..
rm -rf /tmp/testcase1.log
popd
}

# Test cases performs credit card application status rejected if Creditscore < 750
function testcase2 {
pushd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard
go build
./creditcard > /tmp/testcase2.log 2>&1 &
pId=$!

response=$(curl -X PUT http://localhost:7777/newaccount -H 'Content-Type: application/json' -d '{"Name":"Sam4","Age":"26","Income":"50100","Address":"SFO","Id":"5"}' --write-out '%{http_code}' --silent --output /dev/null)
response1=$(curl -X PUT http://localhost:7777/credit -H 'Content-Type: application/json' -d '{"Id":"5","creditScore":"650"}' --write-out '%{http_code}' --silent --output /dev/null)

kill -9 $pId
if [ $response -eq 200  ] && [ $response1 -eq 200  ] && [[ "echo $(cat /tmp/testcase2.log)" =~ "c" ]]
    then 
        echo "PASS"
    else
        echo "FAIL"
fi
cd ..
rm -rf /tmp/testcase2.log
popd
}


# Test cases performs invalid applicant when age address or income data is not matching requirements
function testcase3 {
pushd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/creditcard
go build
./creditcard > /tmp/testcase3.log 2>&1 &
pId=$!

response=$(curl -X PUT http://localhost:7777/newaccount -H 'Content-Type: application/json' -d '{"Name":"Sam4","Age":"26","Income":"5010","Address":"SFO","Id":"6"}' --write-out '%{http_code}' --silent --output /dev/null)

kill -9 $pId
if [ $response -eq 200  ]  && [[ "echo $(cat /tmp/testcase3.log)" =~ "Applicant is not eligible to apply for creditcard" ]]
    then 
        echo "PASS"
    else
        echo "FAIL"
fi
cd ..
rm -rf /tmp/testcase3.log
popd
}
