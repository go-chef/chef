# Inspec tests for the required_recipe chef api go module
#

describe command('/go/src/testapi/bin/required_recipe') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(%r{List required_recipe: file '/tmp/required'}) }
end
