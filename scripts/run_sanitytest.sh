#!/bin/bash

RULESPATH=$GOPATH/src/github.com/project-flogo/rules
export FILENAME="RulesSanityReport.html"
HTML="<!DOCTYPE html>
<html><head><style>table {font-family: arial, sans-serif;border-collapse: collapse;margin: auto;}td,th {border: 1px solid #dddddd;text-align: left;padding: 8px;}th {background: #003399;text-align: center;color: #fff;}body {padding-right: 15px;padding-left: 15px;margin-right: auto;margin-left: auto;}label {font-weight: bold;}.test-report h1 {color: #003399;}.summary,.test-report {text-align: center;}.success {background-color: #79d279;}.error {background-color: #ff3300;}.summary-tbl {font-weight: bold;}.summary-tbl td {border: none;}</style></head><body>    <section class=test-report><h1>Rules Sanity Report</h1></section><section class=summary><h2>Summary</h2><table class="summary-tbl"><tr><td>Number of test cases passed </td> <td> </td></tr><tr><td>Number of test cases failed </td> <td> </td></tr><td>Total test cases</td><td> </td></tr></tr></table></section><section class=test-report><table><tr><th>Recipe</th><th> Testcase </th><th>Status</th><tr></tr> </table></html>"

echo $HTML >> $RULESPATH/scripts/$FILENAME
PASS_COUNT=0 FAIL_COUNT=0

# Fetch list of sanity.sh files in examples folder
function get_sanitylist()
{
    cd $RULESPATH/examples
    find | grep sanity.sh > file.txt
    readarray -t array < file.txt
    for EXAMPLE in "${array[@]}"
    do
        echo "$EXAMPLE"
        RECIPE=$(echo $EXAMPLE  | sed -e 's/\/sanity.sh//g' | sed -e 's/\.\///g' | sed -e 's/\//-/g')
        execute_testcase
    done
}

# Execute and obtain testcase status (pass/fail)
function execute_testcase()
{
    echo $RECIPE
    source $EXAMPLE
    TESTCASE_LIST=($(get_test_cases))
    sleep 10       
    for ((i=0;i < ${#TESTCASE_LIST[@]};i++))
    do
        TESTCASE=$(${TESTCASE_LIST[i]})
        sleep 10
        if [[ $TESTCASE == *"PASS"* ]];  then
            echo "$RECIPE":"Passed"
            PASS_COUNT=$((PASS_COUNT+1))
            sed -i "s/<\/tr> <\/table>/<tr><td>$RECIPE<\/td><td>${TESTCASE_LIST[i]}<\/td><td  class="success">PASS<\/td><\/tr><\/tr> <\/table>/g" $RULESPATH/scripts/$FILENAME
        else
            echo "$RECIPE":"Failed"
            FAIL_COUNT=$((FAIL_COUNT+1))
            sed -i "s/<\/tr> <\/table>/<tr><td>$RECIPE<\/td><td>${TESTCASE_LIST[i]}<\/td><td  class="error">FAIL<\/td><\/tr><\/tr> <\/table>/g" $RULESPATH/scripts/$FILENAME
        fi
    done
}

get_sanitylist

# Update testcase count in html report
sed -i s/"passed <\/td> <td>"/"passed <\/td> <td>$PASS_COUNT"/g $RULESPATH/scripts/$FILENAME
sed -i s/"failed <\/td> <td>"/"failed <\/td> <td>$FAIL_COUNT"/g $RULESPATH/scripts/$FILENAME
sed -i s/"cases<\/td><td>"/"cases<\/td><td>$((PASS_COUNT+FAIL_COUNT))"/g $RULESPATH/scripts/$FILENAME