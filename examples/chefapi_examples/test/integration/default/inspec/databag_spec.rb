# Inspec tests for the databag chef api go module
#
describe command('/go/src/chefapi_test/bin/databag') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/^List data bags before/) }
  its('stdout') { should match(%r{^Data bag created &\{URI:https://localhost/organizations/test/data/testbag\}}) }
  its('stdout') { should match(%r{^List data bags after create testbag => https://localhost/organizations/test/data/testbag}) }
  its('stdout') { should match(/^Created item map\[(?=.*id:item1\s*)(?=.*foo:test123\s*).*\]/) }
  its('stdout') { should match(%r{^List bag items item1 => https://localhost/organizations/test/data/testbag/item1}) }
  its('stdout') { should match(/^Initial item map\[(?=.*id:item1\s*)(?=.*foo:test123\s*).*\]/) }
  its('stdout') { should match(/^Update item map\[(?=.*id:item1\s*)(?=.*foo:update123\s*).*\]/) }
  its('stdout') { should match(/^Updated item map\[(?=.*id:item1\s*)(?=.*foo:update123\s*).*\]/) }
  its('stdout') { should match(/^Deleted item/) }
  its('stdout') { should match(/^List bag items after delete/) }
  its('stdout') { should match(/^Data bag deleted &\{Name:testbag JsonClass:Chef::DataBag ChefType:data_bag\}/) }
  its('stdout') { should match(/^List data bags after delete/) }
end
