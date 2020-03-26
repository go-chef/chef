# Inspec tests for the search chef api go module
#

describe command('/go/src/chefapi_test/bin/search_pagination') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/^List nodes from Exec query Total:1200/) }
  its('stdout') { should match(/^List nodes from partial search Total:1200/) }
end
