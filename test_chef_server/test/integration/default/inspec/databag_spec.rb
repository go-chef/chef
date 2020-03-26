# Inspec tests for the databag chef api go module
#

describe command('/go/src/chefapi_test/bin/databag') do
  its('stderr') { should match(%r{^Issue recreating databag1. POST https://testhost/organizations/test/data: 409}) }
  its('stderr') { should match(%r{^Issue getting nothere. GET https://testhost/organizations/test/data/nothere: 404}) }
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/^List initial databags\s*$/) }
  its('stdout') { should match(%r{^Added databag1 \&\{https://testhost/organizations/test/data/databag1\}}) }
  its('stdout') { should match(%r{^List databags after adding databag1 databag1 => https://testhost/organizations/test/data/databag1}) }
  its('stdout') { should match(/^Create databag1::item1 \<nil\>/) }
  its('stdout') { should match(/^Update databag1::item1 \<nil\>/) }
  its('stdout') { should match(%r{^List databag1 items item1 => https://testhost/organizations/test/data/databag1/item1}) }
  its('stdout') { should match(/^Get databag1::item1 map\[id:item1 type:password value:next\]/) }
  its('stdout') { should match(/^Delete databag1::item1 \<nil\>/) }
  its('stdout') { should match(/^List databag1 items after delete/) }
  its('stdout') { should match(/^Delete databag1 \&\{Name:databag1 JsonClass:Chef::DataBag ChefType:data_bag/) }
  its('stdout') { should match(/^List databags after cleanup\s*$/) }
end
