# Decision Table Usage

This example demonstrates usage decision table with student analysis example.

### Pre-requisites
* Go 1.11
* The **GOPATH** environment variable on your system must be set properly

## Setup and Usage

The following rules are used in the example:

* `studentcare`: Gets fired when `student.careRequired` is true and invokes the `printstudentinfo` function. This function prints the Student name and comments. 
* `studentanalysis`: Gets fired when `studentanalysis.name` and `student.name` are same. It invokes  `dtableservice` which analyses the student and updates the `student.careRequired` and `student.comments` accordingly. 

### Steps

```sh
cd $GOPATH/src/github.com/project-flogo/rules/examples/dtable
go run main.go
```

### Testing

s1 and s2 student tuples are Saved. se1 and se2 are events to analyse s1 and s2 students against the `dtable-file.xlsx` file respectively.
You should see following outputs in the logs:
```
Student Name:  s1  Comments:  “additional study hours required”
```
```
Student Name:  s2  Comments:  “little care can be taken to achieve grade-a”
```