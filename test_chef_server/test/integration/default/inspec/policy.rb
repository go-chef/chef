# Inspec tests for the policy chef api go module
#

describe command('/go/src/github.com/go-chef/chef/testapi/bin/policy') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stderr') { should match(/Issue getting nothere.*404/) }
  its('stdout') { should match(/List policies map\[testsamp:\{Uri/) }
  its('stdout') { should match(/Get testsamp.* map\[revisions:map/) }
  its('stdout') { should match(/Get testsamp.* revision \{RevisionID/) }
  its('stdout') { should match(/Delete revision .*from testsamp.*\{RevisionID/) }
end
