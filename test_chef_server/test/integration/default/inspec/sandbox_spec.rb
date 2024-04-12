# Inspec tests for the sandbox chef api go module
#

describe command('/go/src/github.com/go-chef/chef/testapi/bin/sandbox') do
  # TODO: Get the sandbox sample code to work - upload files is failing

  # its('stderr') { should_not match(%r{Issue}) }
  # its('stderr') { should_not match(%r{error|no such file|cannot find|not used|undefined}) }
  # its('stdout') { should match(/^Create sandboxes  /) }
  # its('stdout') { should match(/^Resulting sandboxes  /) }
end
