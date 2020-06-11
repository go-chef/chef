# Inspec tests for the organization chef api go module
#

<<<<<<< HEAD:test_chef_server/test/integration/default/inspec/organization_spec.rb
describe command('/go/src/testapi/bin/organization') do
  its('stderr') { should match(%r{^Issue creating org: {org1 organization1 } POST https://testhost/organizations: 409$}) }
  its('stderr') { should match(/^Issue creating org: {org1 organization1 } 409$/) }
=======
describe command('/go/src/chefapi_test/bin/organization') do
  its('stderr') { should match(%r{^Issue creating org: {org1 organization1 } POST https://testhost/organizations: 409$}) }
>>>>>>> master:examples/chefapi_examples/test/integration/default/inspec/organization_spec.rb
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/^List initial organizations map\[test.*test\]$/) }
  its('stdout') { should match(/^Added org1 {org1-validator -----BEGIN RSA PRIVATE KEY-----/) }
  its('stdout') { should match(/^Added org1 again {  }$/) }
  its('stdout') { should match(/^Added org2 {org2-validator -----BEGIN RSA PRIVATE KEY-----.*$/) }
  its('stdout') { should match(/^Get org1 {org1 organization1 [0-9a-f]+}$/) }
  its('stdout') { should match(%r{^List organizations After adding org1 and org2 map(?=.*org2:https://testhost/organizations/org2)(?=.*test:https://testhost/organizations/test)(?=.*org1:https://testhost/organizations/org1)}) }
  its('stdout') { should match(/^Update org1 {org1 new_organization1 }/) }
  its('stdout') { should match(/^Get org1 after update {org1 new_organization1 [0-9a-f]+}/) }
  its('stdout') { should match(/^Delete org2 <nil>/) }
  its('stdout') { should match(%r{^List organizations after deleting org2 map(?=.*test:https://testhost/organizations/test)(?=.*org1:https://testhost/organizations/org1)}) }
  its('stdout') { should match(/^Result from deleting org1 <nil>/) }
  its('stdout') { should match(%r{^List organizations after cleanup map\[test:https://testhost/organizations/test\]}) }
end
