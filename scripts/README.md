## Sanity Testing

* There is a shell script file `run_sanitytest.sh` that performs sanity testing against rules/examples and generates html report.

* This script file checks for all available `sanity.sh` files inside rules/examples and run tests against individual `sanity.sh` file 


* To run sanity tests

```
cd $GOPATH/src/github.com/project-flogo/rules/scripts
./run_sanitytest.sh
```

* Testcase status of each example is updated in the html report and test report is made available in scripts folder.


### Contributing

If you're adding a new rules example, optionally you can add sanity test file with name `sanity.sh`. Below is the template used for creating test file.

```
#!/bin/bash

function get_test_cases {
    local my_list=( testcase1 )
    echo "${my_list[@]}"
}

function testcase1 {
# Add detailed steps to execute the test case
}    
```
Sample sanity test file can be found at 
```
$GOPATH/src/github.com/project-flogo/rules/examples/flogo/simple/sanity.sh
```