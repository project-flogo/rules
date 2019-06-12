## Sanity Test Framework

* This is a shell script based test framework that performs sanity testing against rules/examples and generates html report.

* This framework will get the list of all sanity.sh files in examples and run tests against each sanity.sh file 


* To execute sanity tests locally , run below commands

```
cd $GOPATH/src/github.com/project-flogo/rules/scripts
./run_sanitytest.sh
```

* Test status of each example is updated in the html report. This test report is made available in GOPATH


### Contributing

If you're adding a new rules example, optionally you can add sanity test file with name 'sanity.sh'. Below is the template used for creating test file.

```
#!/bin/bash

function get_test_cases {
    local my_list=( testcase1 )
    echo "${my_list[@]}"
}

function testcase1 {

}    
```

* Inside Testcase function user need to add detailed steps to execute the test case.
* In order to execute all the test cases, Testcase functions created needs to be added to my_list under get_test_cases function.