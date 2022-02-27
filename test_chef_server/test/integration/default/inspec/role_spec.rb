# Inspec tests for the role chef api go module
#

describe command('/go/src/testapi/bin/role') do
  its('stderr') { should match(%r{^Issue recreating role1: POST https://testhost/organizations/test/roles: 409}) }
  its('stderr') { should match(/^Issue recreating role1: 409/) }
  its('stderr') { should match(%r{^Issue getting nothere: GET https://testhost/organizations/test/roles/nothere: 404}) }
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(%r{^Added role1 uri => https://testhost/organizations/test/roles/role1}) }
  its('stdout') { should match(%r{^Added roleNR uri => https://testhost/organizations/test/roles/roleNR}) }
  its('stdout') { should match(%r{^List roles after adding role1 role1 => https://testhost/organizations/test/roles/role1}) }
  its('stdout') { should match(/^Get role1 &\{Name:role1 ChefType:role DefaultAttributes:map\[(?=.*git_repo:here.git)(?=.*users:\[root moe\]).*\] Description:Test role EnvRunList:map\[(?=.*en2:\[recipe\[foo2\]\]).*(?=.*en1:\[recipe\[foo1\] recipe\[foo2\]\]).*JsonClass:Chef::Role OverrideAttributes:map\[env:map\[(?=.*mine:ample)(?=.*yours:full).*\]\] RunList:\[recipe\[foo\] recipe\[baz\] role\[banana\]\]\}/) }
  its('stdout') { should match(/^Update role1 &\{Name:role1 ChefType:role DefaultAttributes:map\[(?=.*git_repo:here.git)(?=.*users:\[root moe\]).*\] Description:Changed Role EnvRunList:map\[(?=.*en1:\[recipe\[foo1\] recipe\[foo2\]\]).*(?=.*en2:\[recipe\[foo2\]\]).*JsonClass:Chef::Role OverrideAttributes:map\[env:map\[(?=.*mine:ample)(?=.*yours:full).*\]\] RunList:\[recipe\[foo\] recipe\[baz\] role\[banana\]\]\}/) }
  its('stdout') { should match(/^Environments for role1 \[_default en1 en2\]/) }
  its('stdout') { should match(/^Environments for role1 map\[run_list:\[recipe\[foo1\] recipe\[foo2\]\]\]/) }
  its('stdout') { should match(/^Delete role1 <nil>/) }
  its('stdout') { should match(/^List roles after cleanup.*/) }
end
