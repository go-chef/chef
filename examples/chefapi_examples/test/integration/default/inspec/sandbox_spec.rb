# Inspec tests for the sandbox chef api go module
#

describe command('/go/src/chefapi_test/bin/sandbox') do
  # TODO: Get the sandbox sample code to work - upload files is failing

  # its('stderr') { should_not match(%r{Issue}) }
  # its('stderr') { should_not match(%r{error}) }
  # its('stderr') { should_not match(%r{no such file}) }
  # its('stderr') { should_not match(%r{cannot find}) }
  # its('stderr') { should_not match(%r{not used|undefined}) }
  # its('stdout') { should match(/^Create sandboxes  /) }
  # its('stdout') { should match(/^Resulting sandboxes  /) }
end
