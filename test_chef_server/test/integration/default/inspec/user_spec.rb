# Inspec tests for the user chef api go module
#

describe command('/go/src/github.com/go-chef/chef/testapi/bin/user') do
  its('stderr') { should match(%r{^Issue creating user POST https://testhost/users: 409}) }
  its('stderr') { should match(%r{^Issue creating user err: POST https://testhost/users: 409}) }
  its('stderr') { should match(/^Issue creating user code: 409/) }
  its('stderr') { should match(/^Issue creating user method: POST/) }
  its('stderr') { should match(%r{^Issue creating user url: https://testhost/users}) }
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(%r{^List initial users map\[(?=.*pivotal:https://testhost/users/pivotal).*\] EndInitialList}) }
  # might want a multi line match here to test for expirationdate, key uri and privatekey
  its('stdout') { should match(%r{^Add usr1 \{Uri:https://testhost/users/usr1 ChefKey:\{Name:default PublicKey:-----BEGIN}) }
  its('stdout') { should match(%r{^Add usr2 \{Uri:https://testhost/users/usr2 ChefKey:\{Name:default PublicKey:-----BEGIN}) }
  its('stdout') { should match(%r{^Add usr3 \{Uri:https://testhost/users/usr3 ChefKey:\{Name: PublicKey: ExpirationDate: Uri: PrivateKey:\}\}}) }
  its('stdout') { should match(%r{^Filter users map\[usr1:https://testhost/users/usr1\]}) }
  its('stdout') { should match(/^Verbose out map\[(?=.*pivotal:)/) }
  its('stdout') { should match(/^Get usr1 \{(?=.*UserName:usr1)(?=.*DisplayName:User1 Fullname)(?=.*Email:user1@domain.io)(?=.*ExternalAuthenticationUid:)(?=.*FirstName:user1)(?=.*LastName:fullname)(?=.*MiddleName:)(?=.*Password:)(?=.*PublicKey:)(?=.*RecoveryAuthenticationEnabled:false).*/) }
  its('stdout') { should match(/^Pivotal user \{(?=.*UserName:pivotal)(?=.*DisplayName:Chef Server Superuser)(?=.*Email:root@localhost.localdomain)(?=.*ExternalAuthenticationUid:)(?=.*FirstName:Chef)(?=.*LastName:Server)(?=.*MiddleName:)(?=.*Password:)(?=.*PublicKey:)/) }
  its('stdout') { should match(%r{^List after adding map\[(?=.*pivotal:https://testhost/users/pivotal)(?=.*usr1:https://testhost/users/usr1).*\] EndAddList}) }
  its('stdout') { should match(/^Get usr1 \{(?=.*UserName:usr1)(?=.*DisplayName:User1 Fullname)(?=.*Email:user1@domain.io)(?=.*ExternalAuthenticationUid:)(?=.*FirstName:user1)(?=.*LastName:fullname)(?=.*MiddleName:)(?=.*Password:)(?=.*PublicKey:)/) }
  its('stdout') { should match(%r{^List after adding map\[(?=.*pivotal:https://testhost/users/pivotal)(?=.*usr1:https://testhost/users/usr1).*\] EndAddList}) }
  # TODO: - update and create new private key
  # TODO - is admin a thing
  its('stdout') { should match(%r{^Update usr1 partial update \{Uri:https://testhost/users/usr1 ChefKey:\{}) }
  its('stdout') { should match(/^Get usr1 after partial update \{(UserName:usr1)(?=.*DisplayName:usr1)(?=.*Email:myuser@samp.com)(?=.*ExternalAuthenticationUid:)(?=.*FirstName:user1)(?=.*LastName:fullname)(?=.*MiddleName:)(?=.*Password:)(?=.*PublicKey:)(?=.*RecoveryAuthenticationEnabled:false).*\}/) }
  its('stdout') { should match(%r{^Update usr1 full update \{Uri:https://testhost/users/usr1 ChefKey:\{Name: PublicKey: ExpirationDate: Uri: PrivateKey:\}}) }
  its('stdout') { should match(/^Get usr1 after full update \{(UserName:usr1)(?=.*DisplayName:usr1)(?=.*Email:myuser@samp.com)(?=.*ExternalAuthenticationUid:)(?=.*FirstName:user)(?=.*LastName:name)(?=.*MiddleName:mid)(?=.*Password:)(?=.*PublicKey:)(?=.*RecoveryAuthenticationEnabled:false).*\}/) }
  its('stdout') { should match(%r{^Update usr1 rename \{Uri:https://testhost/users/usr1new ChefKey:\{.*\}}) }
  its('stdout') { should match(/^Get usr1 after rename \{(UserName:usr1new)(?=.*DisplayName:usr1)(?=.*Email:myuser@samp.com)(?=.*ExternalAuthenticationUid:)(?=.*FirstName:user)(?=.*LastName:name)(?=.*MiddleName:mid Password:)(?=.*PublicKey:)(?=.*RecoveryAuthenticationEnabled:false).*\}/) }
  its('stdout') { should match(%r{^Delete usr1 DELETE https://testhost/users/usr1: 404}) }
  its('stdout') { should match(/^Delete usr1new <nil>/) }
  its('stdout') { should match(%r{^List after cleanup map\[(?=.*pivotal:https://testhost/users/pivotal).*\] EndCleanupList}) }
end
