# Inspec tests for the stats chef api go module
#

describe command('/go/src/testapi/bin/stats') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/List stats json format: \[map/) }
end
