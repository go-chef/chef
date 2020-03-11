# Inspec tests for the user chef api go module
#

describe command('/go/src/chefapi_test/bin/userkey') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(%r{^List keys }) }
end
