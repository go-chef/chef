# Inspec tests for the cookbook chef api go module
#
describe command('/go/src/github.com/go-chef/chef/testapi/bin/cookbook') do
  its('stderr') do
    should match(%r{^Issue getting cookbook nothere: GET https://testhost/organizations/test/cookbooks/nothere: 404})
  end
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stderr') { should_not match(/testbook/) }
  its('stderr') { should_not match(/sampbook/) }
  its('stderr') { should_not match(/Issue loading/) }
  its('stdout') do
    should match(%r{^List initial cookbooks (?=.*sampbook => https://testhost/organizations/test/cookbooks/sampbook\n\s*\* 0.2.0)(?=.*testbook => https://testhost/organizations/test/cookbooks/testbook\n\s*\* 0.2.0).*EndInitialList}m)
  end
  # output from get cookbook is odd
  its('stdout') { should match(/^Get cookbook testbook/) }
  its('stdout') do
    should match(%r{^Get cookbook versions testbook testbook => https://testhost/organizations/test/cookbooks/testbook\n\s*\* 0.2.0\n\s*\* 0.1.0}m)
  end
  its('stdout') do
    should match(%r{^Get cookbook versions sampbook sampbook => https://testhost/organizations/test/cookbooks/sampbook\n\s*\* 0.2.0\n\s*\* 0.1.0}m)
  end
  its('stdout') { should match(/^Get specific cookbook version testbook {CookbookName:testbook/) }
  its('stdout') { should match(/^Get all recipes \[sampbook testbook\]/) }
  its('stdout') { should match(/^Delete testbook 0.1.0 <nil>/) }
  its('stdout') do
    should match(%r{^Final cookbook list sampbook => https://testhost/organizations/test/cookbooks/sampbook\n\s*\* 0.2.0}m)
  end
  its('stdout') do
    should match(%r{^Final cookbook versions sampbook sampbook => https://testhost/organizations/test/cookbooks/sampbook\n\s*\* 0.2.0\n\s*\* 0.1.0}m)
  end
end
