# Inspec tests for the client chef api go module
#

describe command('/go/src/testapi/bin/client') do
  its('stderr') { should match(%r{^Couldn't recreate client client1.  POST https://testhost/organizations/test/clients: 409}) }
  its('stderr') { should match(/^Couldn't recreate client client1.  Code 409/) }
  its('stderr') { should match(/^Couldn't recreate client client1.  Msg Client already exists/) }
  its('stderr') { should match(/^Couldn't recreate client client1.  Text \{"error":\["Client already exists"\]\}/) }
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(%r{^List initial clients test-validator => https://testhost/organizations/test/clients/test-validator}) }
  its('stdout') { should match(/^Define client1 \{Name:client1 ClientName: Validator:false Admin:false CreateKey:true\}/) }
  its('stdout') { should match(%r{^Added client1 &\{Uri:https://testhost/organizations/test/clients/client1 ChefKey:\{Name:default PublicKey:-----BEGIN PUBLIC KEY-----MIIB.*ExpirationDate:infinity Uri:.*PrivateKey:-----BEGIN RSA PRIVATE KEY}) }
  its('stdout') { should match(%r{^Added client2 &\{Uri:https://testhost/organizations/test/clients/client2 ChefKey:\{Name: PublicKey: ExpirationDate: Uri: PrivateKey:\}\}}) }
  # TODO: are OrgName and Uri part of the get output
  its('stdout') { should match(/Get client1 \{Name:client1 ClientName:client1 OrgName:test Validator:false JsonClass:Chef::ApiClient ChefType:client\}/) }
  its('stdout') { should match(/Get client2 \{Name:client2 ClientName:client2 OrgName:test Validator:true JsonClass:Chef::ApiClient ChefType:client\}/) }
  # TODO: are orgname and uri part of the output
  its('stdout') { should match(/Update client1 &\{Name:clientchanged ClientName:clientchanged OrgName: Validator:true JsonClass:Chef::ApiClient ChefType:client\}/) }
  its('stdout') { should match(/Update client2 &\{Name:client2 ClientName:client2 OrgName: Validator:true JsonClass:Chef::ApiClient ChefType:client\}/) }
  its('stdout') { should match(/^Delete client1 <nil>/) }
  its('stdout') { should match(/^Delete client2 <nil>/) }
  its('stdout') { should match(%r{^List clients after cleanup test-validator => https://testhost/organizations/test/clients/test-validator}) }
end
