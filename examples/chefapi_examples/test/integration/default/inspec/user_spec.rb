# Inspec tests for the user chef api go module
#

describe command('/go/src/chefapi_test/bin/user') do
  its('stderr') { should match(%r{^Issue creating user: POST https://localhost/users: 409}) }
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(%r{^List initial users map\[(?=.*pivotal:https://localhost/users/pivotal).*\] EndInitialList}) }
  its('stdout') { should match(%r{^Add usr1 \{https://localhost/users/usr1 -----BEGIN RSA}) }
  its('stdout') { should match(%r{^Filter users map\[usr1:https://localhost/users/usr1\]}) }
  its('stdout') { should match(/^Verbose out map\[(?=.*pivotal:)/) }
  its('stdout') { should match(/^Get usr1 \{UserName:usr1 DisplayName:User1 Fullname Email:user1@domain.io ExternalAuthenticationUid: FirstName:user1 FullName: LastName:fullname MiddleName: Password: PublicKey:-----BEGIN/) }
  its('stdout') { should match(/^Pivotal user \{UserName:pivotal DisplayName:Chef Server Superuser Email:root@localhost.localdomain ExternalAuthenticationUid: FirstName:Chef FullName: LastName:Server MiddleName: Password: PublicKey:-----BEGIN/) }
  its('stdout') { should match(%r{^List after adding map\[(?=.*pivotal:https://localhost/users/pivotal)(?=.*usr1:https://localhost/users/usr1).*\] EndAddList}) }
  # its('stdout') { should match(/^Update usr1     /) }
  # its('stdout') { should match(/^Get usr1 after update        /) }
  its('stdout') { should match(/^Delete usr1 <nil>/) }
  its('stdout') { should match(%r{^List after cleanup map\[(?=.*pivotal:https://localhost/users/pivotal).*\] EndCleanupList}) }
end
