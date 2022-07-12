# Inspec tests for the search chef api go module
#

describe command('/go/src/testapi/bin/search') do
  its('stderr') { should match(/^Issue building invalid query statement is malformed/) }
  its('stderr') { should_not match(/node/) }
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(%r{^List indexes map\[(?=.*node:https://testhost/organizations/test/search/node)(?=.*role:https://testhost/organizations/test/search/role)(?=.*client:https://testhost/organizations/test/search/client)(?=.*environment:https://testhost/organizations/test/search/environment).*\] EndIndex}) }
  its('stdout') { should match(/^List new query node\?q=name:node\*\&rows=2\&sort=X_CHEF_id_CHEF_X asc\&start=0/) }
  its('stdout') { should match(/^List nodes from query \{Total:2 Start:0 Rows:\[/) }
  its('stdout') { should match(/^List nodes from Exec query \{Total:1 Start:0 Rows:\[/) }
  its('stdout') { should match(/^List nodes from all nodes Exec query \{Total:4 Start:0 Rows:\[/) }
  its('stdout') { should match(/^List nodes from partial search \{Total:4 Start:0 Rows:\[/) }
  its('stdout') { should match(/^List nodes from query JSON format \{Total:2 Start:0 Rows:\[/) }
  its('stdout') { should match(/^List 2nd set of nodes from query JSON format \{Total:2 Start:0 Rows:\[/) }
  its('stdout') { should match(/^List nodes from Exec query JSON format \{Total:1 Start:0 Rows:\[/) }
  its('stdout') { should match(/^List nodes from all nodes Exec query JSON format \{Total:4 Start:0 Rows:\[/) }
  its('stdout') { should match(/^Partial search JSON format row: 0 rawjson: {"name":"node0"}/) }
end
