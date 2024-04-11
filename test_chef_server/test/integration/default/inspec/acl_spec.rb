# Inspec tests for the ACL chef api go module
#

describe command('/go/src/github.com/go-chef/chef/testapi/bin/acl') do
  its('stderr') { should_not match(/Issue/) }
end
