I'll add some doc about testing and development. I'll take the request as a really good suggestion.
Thanks

Basic ideas for adding an endpoint, assume you will repeat this process and each step multiple times:

Design the data structures for a new end point

Create go tests for each func.

Write the go code.

Add examples to exercise the functions in the chefapi_examples cookbook.

Create go code to call the functions. See the files in go-chef/chef/testapi/cmd. Many times extra setup is required.
Create a shell script to call the go code. See the files in go-chef/chef/testapi/bin.
Write Inspec tests to verify the output. Be aware that go map output does not alway have consistent ordering. The matching regex structures are rarely fun to write.
Commit the go-chef/chef (api client go code at least) changes to a branch and push to GitHub. The chefapi.rb recipe would need to be updated to pull from the correct repo.
kitchen converge. - installs the code, spins up a chef server to test against
kitchen login - run the /bin commands and iterate on fixing the code. Both cmd code and chef client code. kitchen converge will install cmd and bin file changes without committing and pushing to GitHub. Changes to the chef api client code need to be committed, pushed and kitchen converge run again to get the changes installed.
While logged in to the virtual box image try adding "-tags" debug the bin command. Example
go run -tags debug ${BASE}/../cmd/client/clients.go ${CHEFUSER} ${KEYFILE} ${CHEFORGANIZATIONURL} ${SSLBYPASS} will produce output that shows the real request body returned from the chef server. Note what you passed in and what you got back. Send updates for the documentation to https://github.com/chef/chef-web-docs
After things work ok. Create an inspec test for the endpoint that calls the bin command for the endpoint. Run kitchen verify to run the full set of integration tests.
Run cookstyle -a in the go-chef/chef/test_chef_server directory

Run go fmt in the go-chef/chef directory

See if any documentation needs to be updated, check to see if any of the scattered TODOs are now now.

Commit and PR
