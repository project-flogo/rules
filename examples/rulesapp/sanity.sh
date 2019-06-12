#!/bin/bash

function get_test_cases {
    local my_list=( testcase1 )
    echo "${my_list[@]}"
}

# This testcase checks for name bob
function testcase1 {
pushd $GOPATH/src/github.com/project-flogo/rules/examples/rulesapp
rm -rf /tmp/testcase1.log
go run main.go > /tmp/testcase1.log 2>&1

if [[ "echo $(cat /tmp/testcase1.log)" =~ "Rule fired" ]]
    then 
        echo "PASS"
    else
        echo "FAIL"
fi
popd
}