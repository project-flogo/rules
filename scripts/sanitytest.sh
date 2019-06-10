#!/bin/bash

RULESPATH=$GOPATH/src/github.com/project-flogo/rules
export FILENAME="RulesSanityReport.html"
HTML="<!DOCTYPE html>
<html><head><style>table {font-family: arial, sans-serif;border-collapse: collapse;margin: auto;}td,th {border: 1px solid #dddddd;text-align: left;padding: 8px;}th {background: #003399;text-align: center;color: #fff;}body {padding-right: 15px;padding-left: 15px;margin-right: auto;margin-left: auto;}label {font-weight: bold;}.test-report h1 {color: #003399;}.summary,.test-report {text-align: center;}.success {background-color: #79d279;}.error {background-color: #ff3300;}.summary-tbl {font-weight: bold;}.summary-tbl td {border: none;}</style></head><body>    <section class=test-report><h1>Rules Sanity Report</h1></section><section class=summary><h2>Summary</h2><table class="summary-tbl"><tr><td>Number of test cases passed </td> <td> </td></tr><tr><td>Number of test cases failed </td> <td> </td></tr><td>Total test cases</td><td> </td></tr></tr></table></section><section class=test-report><table><tr><th>Recipe</th><th> Testcase </th><th>Status</th><tr></tr> </table></html>"

echo $HTML >> $GOPATH/$FILENAME
p=0 q=0 r=0

# Fetch list of sanity.sh files in examples folder
function get_list()
{
    cd $RULESPATH/examples
    find | grep sanity.sh > file.txt
    readarray -t array < file.txt
    for EXAMPLE in "${array[@]}"
    do
        echo "$EXAMPLE"
        RECIPE=${EXAMPLE%/sanity.sh*}
        RECIPE=${RECIPE##*./}
        RECIPE=$(echo $RECIPE | sed -e 's/\//-/g')
        testcase_status
    done
}

# Execute and obtain testcase status (pass/fail)
function testcase_status()
{
    echo $RECIPE
    source $EXAMPLE
    value=($(get_test_cases))
    sleep 10       
    for ((i=0;i < ${#value[@]};i++))
    do
        value1=$(${value[i]})
        sleep 10
        if [[ $value1 == *"PASS"* ]];  then
            echo "$RECIPE":"Passed"
            q=$((q+1))
            sed -i "s/<\/tr> <\/table>/<tr><td>$RECIPE<\/td><td>${value[i]}<\/td><td  class="success">PASS<\/td><\/tr><\/tr> <\/table>/g" $GOPATH/$FILENAME
        else
            echo "$RECIPE":"Failed"
            r=$((r+1))
            sed -i "s/<\/tr> <\/table>/<tr><td>$RECIPE<\/td><td>${value[i]}<\/td><td  class="error">FAIL<\/td><\/tr><\/tr> <\/table>/g" $GOPATH/$FILENAME
        fi
        p=$((p+1))
    done
}

get_list

# Update testcase count in html report
sed -i s/"passed <\/td> <td>"/"passed <\/td> <td>$q"/g $GOPATH/$FILENAME
sed -i s/"failed <\/td> <td>"/"failed <\/td> <td>$r"/g $GOPATH/$FILENAME
sed -i s/"cases<\/td><td>"/"cases<\/td><td>$p"/g $GOPATH/$FILENAME