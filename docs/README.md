# Flogo Rules
1. [Introduction](#introduction)
2. [Terminology](#Terminology)

## Introduction

![Rule High Level Architecture](https://github.com/project-flogo/rules/blob/feature-improve-docs/docs/hl-architecture.jpeg)

## Terminology
### trigger
`trigger` is a flogo/trigger which receives data from external sources. `rest`, `kafka` & `timer` are the few `triggers` to name here. For comprehensive list you can visit [project-flogo/contrib](https://github.com/project-flogo/contrib)
### activity
`activity` is a flogo/activity which implements common application logic in reusalble mannaer. `rest`, `kafka`, `sqlquery`, `log` are the few `activities` to name here. For comprehensive list you can visit [project-flogo/contrib](https://github.com/project-flogo/contrib)
### ruleaction
ruleaction is a flogo/action which implements typical rule engine capabilities.
### tuple
`tuple` represents an event or a business object and provides runtime data to rule engine's runtime.
### model
`model` consists of various definitions used by rule engines's run time.
### rule
`rule` constitutes of multiple `conditions` and an `action`. The `rule` triggers when all of its conditions pass.
### condition
Here `condition` is referred to as a rule condition. It is an expression involving one or more tuple types. When the expression evaluates to true, the condition passes. In order to optimize a Rule's evaluation, the Rule network needs to know of the `TupleTypes` and the properties of the `TupleType` which participate in the Condition evaluation. These are provided when constructing the condition and adding it to the rule.
### action
Here `action` is referred to as a rule action. It is a `function` or `flogo/activity` or `flogo/action` that is invoked each time when a matching combination of tuples are found that results in all of its `rule conditions` are evaluated to true. 