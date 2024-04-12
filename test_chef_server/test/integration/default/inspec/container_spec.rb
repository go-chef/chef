# Inspec tests for the container chef api go module
#

describe command('/go/src/github.com/go-chef/chef/testapi/bin/container') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/^List initial containers (?=.*clients)(?=.*cookbook_artifacts)/) }
  its('stdout') { should match(%r{^Added container1 \&\{https://testhost/organizations/test/containers/container1\}}) }
  its('stdout') { should match(/^List containers after adding container1 (?=.*clients)(?=.*cookbook_artifacts)(?=.*container1)/) }
  its('stdout') { should match(/^Get container1 \{ContainerName:container1 ContainerPath:container1\}/) }
  its('stdout') { should match(/^Get environment \{ContainerName:environments ContainerPath:environments\}/) }
  its('stdout') { should match(/^Delete container1 <nil>/) }
end
