# Examples for using the api

## Coded samples
The cmd directory has examples of using the client code to exercise the chef api end points.
The bin directory has commands for invoking the code.
The inspec directory has inspec tests to verify the output from invoking the code.
The chefapi_examples cookbook creates a chef server instance. Test kitchen can run the sample code and verify
the resulting output using inspec tests.

## Cookbook and kitchen testing
Run kitchen converge to create a chef server instance in a linux image running under vagrant.
Run kitchen verify to run the go examples using the client api and to confirm
that the results are as expected.

## Looking at the output from an api call
Sometimes you might want to see what output gets created by an api call.  Inspec tends to hide
and mask the output. You can use kitchen login to access the linux image. Use "cd /go/src/testapi/bin"
to access the bin directory and run any of the commands to run the api sample use code and see
the results. Running the bin commands by adding the -tags debug option will show more detail.

## Creating a client
On the test image /go/src/github.com/go-chef/chef/testapi/testapi.go has code that creates a client
for use by other api calls. For the purposes of testing it uses the pivotal user and key
for all the tests.
