# Inspec tests for the search chef api go module
#

describe command('/go/src/github.com/go-chef/chef/testapi/bin/search_pagination') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/^List nodes from Exec query Total:50 Rows:50/) }
  its('stdout') { should match(/^List nodes detail from Exec query.*(?=.*name:node0)(?=.*name:node49)/) }
  its('stdout') { should match(/^List nodes detail from Exec query.*(?=.*default:map)(?=.*normal:map)/) }
  its('stdout') { should match(/^List nodes from Partial Exec Total:50 Rows:50/) }
  its('stdout') { should match(/^List nodes detail from Partial Exec.*(?=.*data:map\[name:node0\])(?=.*data:map\[name:node49\])/) }
end
