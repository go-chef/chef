# Inspec tests for the search chef api go module
#

describe command('/go/src/chefapi_test/bin/search') do
  its('stderr') { should match(/^Issue building invalid query statement is malformed/) }
  its('stderr') { should_not match(/node/) }
  its('stdout') { should match(%r{^List indexes map\[(?=.*node:https://localhost/organizations/test/search/node)(?=.*role:https://localhost/organizations/test/search/role)(?=.*client:https://localhost/organizations/test/search/client)(?=.*environment:https://localhost/organizations/test/search/environment).*\] EndIndex}) }
  its('stdout') { should match(/^List new query node\?q=name:node\*\&rows=1000\&sort=X_CHEF_id_CHEF_X asc\&start=0/) }
  its('stdout') { should match(/^List nodes from query \{Total:2 Start:0 Rows:\[/) }
  its('stdout') { should match(/^List nodes from Exec query \{Total:1 Start:0 Rows:\[/) }
  its('stdout') { should match(/^List nodes from all nodes Exec query \{Total:4 Start:0 Rows:\[/) }
  its('stdout') { should match(/^List nodes from partial search \{Total:4 Start:0 Rows:\[/) }
end
