# Inspec tests for the authenticate_user chef api go module
#

describe command('/go/src/github.com/go-chef/chef/testapi/bin/authenticate') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/^Authenticate with a valid password <nil>/) }
  its('stdout') do
    should match(%r{^Authenticate with an invalid password POST https://testhost/authenticate_user: 401})
  end
end
