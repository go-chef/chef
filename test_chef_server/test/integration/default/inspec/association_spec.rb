# Inspec tests for the associations chef api go module
#
describe command('/go/src/github.com/go-chef/chef/testapi/bin/association') do
  its('stderr') { should match(/^Issue inviting a user \{User:nouser\} .* 404/) }
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stderr') { should_not match(/testbook/) }
  its('stderr') { should_not match(/sampbook/) }
  its('stdout') { should match(%r{^Invited user \{User:usrinvite\} \{Uri:https://testhost/organizations/test/association_requests/[a-f0-9]+ OrganizationUser:\{UserName:pivotal\} Organization:\{Name:test\} User:\{Email:usrauth@domain.io FirstName:usr\}\}}) }
  its('stdout') { should match(%r{^Invited user \{User:usr2invite\} \{Uri:https://testhost/organizations/test/association_requests/[a-f0-9]+ OrganizationUser:\{UserName:pivotal\} Organization:\{Name:test\} User:\{Email:usr22auth@domain.io FirstName:usr22\}\}}) }
  its('stdout') { should match(/^Invitation id for usr2invite [a-f0-9]+/) }
  its('stdout') { should match(/^Invitation list \[(?=.*\{Id:[a-f0-9]+ UserName:usr2invite\})(?=.*\{Id:[a-f0-9]+ UserName:usrinvite\})/) }
  its('stdout') { should match(/^Deleted invitation [a-f0-9]+ for usrinvite \{Id:[a-f0-9]+ Orgname:test Username:usrinvite\}/) }
  its('stdout') { should match(/^User added: \{Username:usradd\}/) }
  its('stdout') { should match(/^Users list: \[\{User:\{Username:usradd\}\}\]/) }
  its('stdout') { should match(/^User details: \{Username:usradd Email:usradd@domain.io DisplayName:UserAdd Fullname FirstName:usr LastName:add PublicKey:/) }
  its('stdout') { should match(/^User deleted: \{Username:usradd Email:usradd@domain.io DisplayName:UserAdd Fullname FirstName:usr LastName:add PublicKey:/) }
end
