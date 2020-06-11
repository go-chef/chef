# Inspec tests for the user chef api go module
#

describe command('/go/src/testapi/bin/userkeys') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stderr') { should match(%r{Error displaying key detail GET https://testhost/users/usr1/keys/default: 404}) }
<<<<<<< HEAD:test_chef_server/test/integration/default/inspec/userkeys_spec.rb
  its('stdout') { should match(%r{^List initial user usr1 keys \[\{Name:default Uri:https://testhost/users/usr1/keys/default Expired:false\}\]}) }
  its('stdout') { should match(%r{^List initial user usr2 keys \[\{Name:default Uri:https://testhost/users/usr2/keys/default Expired:false\}\]}) }
  its('stdout') { should match(/^List initial user usr3 keys \[\]/) }
  its('stdout') { should match(%r{^Add usr1 key \{Name: Uri:https://testhost/users/usr1/keys/newkey Expired:false\}}) }
  its('stdout') { should match(/^List after add usr1 keys \[\{(?=.*newkey)(?=.*default).*\}\]/) }
  its('stdout') { should match(%r{^Add usr3 key \{Name: Uri:https://testhost/users/usr3/keys/default Expired:false\}}) }
=======
  its('stdout') { should match(%r{^List initial user usr1 keys \[\{KeyName:default Uri:https://testhost/users/usr1/keys/default Expired:false\}\]}) }
  its('stdout') { should match(%r{^List initial user usr2 keys \[\{KeyName:default Uri:https://testhost/users/usr2/keys/default Expired:false\}\]}) }
  its('stdout') { should match(/^List initial user usr3 keys \[\]/) }
  its('stdout') { should match(%r{^Add usr1 key \{KeyName: Uri:https://testhost/users/usr1/keys/newkey Expired:false\}}) }
  its('stdout') { should match(/^List after add usr1 keys \[\{(?=.*newkey)(?=.*default).*\}\]/) }
  its('stdout') { should match(%r{^Add usr3 key \{KeyName: Uri:https://testhost/users/usr3/keys/default Expired:false\}}) }
>>>>>>> master:examples/chefapi_examples/test/integration/default/inspec/userkeys_spec.rb
  its('stdout') { should match(/^List after add usr3 keys \[\{(?=.*default).*\}\]/) }
  its('stdout') { should match(/^Key detail usr1 default \{Name:default/) }
  its('stdout') { should match(/^Key update output usr1 default \{Name:default .*N0AIhUh7Fw1\+gQtR\+.*\}/) }
  its('stdout') { should match(/^Updated key detail usr1 default \{Name:default .*N0AIhUh7Fw1\+gQtR\+.*\}/) }
  its('stdout') { should match(/^List delete result usr1 keys \{Name:default .*N0AIhUh7Fw1\+gQtR\+.*\}/) }
end
