# Inspec tests for the cookbook chef api go module
#
describe command('/go/src/chefapi_test/bin/cookbook') do
  its('stderr') { should match(%r{^Issue getting cookbook nothere: GET https://localhost/organizations/test/cookbooks/nothere: 404}) }
  its('stderr') { should_not match(/testbook/) }
  its('stderr') { should_not match(/sampbook/) }
  its('stderr') { should_not match(/Issue loading/) }
  its('stdout') { should match(%r{^List initial cookbooks (?=.*sampbook => https://localhost/organizations/test/cookbooks/sampbook\n\s*\* 0.2.0)(?=.*testbook => https://localhost/organizations/test/cookbooks/testbook\n\s*\* 0.2.0).*EndInitialList}m) }
  # output from get cookbook is odd
  its('stdout') { should match(/^Get cookbook testbook/) }
  its('stdout') { should match(%r{^Get cookbook versions testbook testbook => https://localhost/organizations/test/cookbooks/testbook\n\s*\* 0.2.0\n\s*\* 0.1.0}m) }
  its('stdout') { should match(%r{^Get cookbook versions sampbook sampbook => https://localhost/organizations/test/cookbooks/sampbook\n\s*\* 0.2.0\n\s*\* 0.1.0}m) }
  its('stdout') { should match(/^Get specific cookbook version testbook {CookbookName:testbook/) }
  its('stdout') { should match(/^Get all recipes \[sampbook testbook\]/) }
  its('stdout') { should match(/^Delete testbook 0.1.0 <nil>/) }
  its('stdout') { should match(%r{^Final cookbook list sampbook => https://localhost/organizations/test/cookbooks/sampbook\n\s*\* 0.2.0}m) }
  its('stdout') { should match(%r{^Final cookbook versions sampbook sampbook => https://localhost/organizations/test/cookbooks/sampbook\n\s*\* 0.2.0\n\s*\* 0.1.0}m) }
end
