# Inspec tests for the http chef api go module
#

describe command('/go/src/github.com/go-chef/chef/testapi/bin/http') do
  its('stderr') { should_not match(/FAILURE/) }
end
