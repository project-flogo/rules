#!/bin/bash

function get_test_cases {
    local my_list=( testcase1 )
    echo "${my_list[@]}"
}

# This Testcase creates flogo rules binary and checks for name bob
function testcase1 {
pushd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/simple
flogo create -f flogo.json
cp functions.go simplerules/src
cd simplerules
flogo build
./bin/simplerules > /tmp/testcase1.log 2>&1 &
pId=$!

response=$(curl --request GET localhost:7777/test/n1?name=Bob --write-out '%{http_code}' --silent --output /dev/null)
response1=$(curl --request GET localhost:7777/test/n2?name=Bob --write-out '%{http_code}' --silent --output /dev/null)

kill -9 $pId
if [ $response -eq 200 ] && [ $response1 -eq 200 ] && [[ "echo $(cat /tmp/testcase1.log)" =~ "Rule fired" ]]
    then 
        echo "PASS"
    else
        echo "FAIL"
fi
cd ..
rm -rf simplerules
popd
}