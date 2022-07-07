# Inspec tests for the policygroup chef api go module
#

describe command('/go/src/testapi/bin/policygroup') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/List policy_groups map.*testgroup:\{Uri/) }
end
