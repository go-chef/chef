# Inspec tests for the node chef api go module
#

describe command('/go/src/chefapi_test/bin/node') do
  its('stderr') { should match(%r{^Couldn't recreate node node1.  POST https://localhost/organizations/test/nodes: 409}) }
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/^List initial nodes map\[\]$/) }
  its('stdout') { should match(/^Define node1 {node1 _default.*Chef::Node \[pwn\]/) }
  its('stdout') { should match(%r{^Added node1 \&\{https://localhost/organizations/test/nodes/node1\}}) }
  its('stdout') { should match(%r{^List nodes after adding node1 map\[node1:https://localhost/organizations/test/nodes/node1\]}) }
  its('stdout') { should match(/^Get node1 {Name:node1 Environment:_default ChefType:node AutomaticAttributes:map\[\] NormalAttributes:map\[\] DefaultAttributes:map\[\] OverrideAttributes:map\[\] JsonClass:Chef::Node RunList:\[recipe\[pwn\]\] PolicyName: PolicyGroup/) }
  its('stdout') { should match(/^Update node1/) }
  its('stdout') { should match(/^Get node1 after update/) }
  its('stdout') { should match(/^Delete node1 <nil>/) }
  its('stdout') { should match(/^List nodes after cleanup map\[\]/) }
end
