# Inspec tests for the cookbook chef api go module
#
describe command('/go/src/github.com/go-chef/chef/testapi/bin/environment') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stderr') { should_not match(/testbook/) }
  its('stderr') { should_not match(/sampbook/) }
  its('stderr') { should_not match(/Issue loading/) }
  its('stdout') { should match(%r{^List initial environments\s*$}) }
end
