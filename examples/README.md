# Examples for using the api
## Coded samples
The code directory has examples of using the client code to exercise the chef api end points.
The bin directory has command for invoking the code.
The chefapi_examples cookbook creates a chef server instance. Then runs the code and verifies
the resulting output using inspec tests.

## Cookbook and kitchen testing
Run kitchen converge to create a chef server instance in a linux image running under vagrant.
Run kitchen verify to code run the go examples that using the client api and to see confirm
that the results are as expected.
