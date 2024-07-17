# Inspec tests for the status chef api go module
#

describe command('/go/src/github.com/go-chef/chef/testapi/bin/status') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') do
    should match(/List status: \{Status:pong Upstreams:map\[(?=.*chef_opensearch:pong)(?=.*chef_sql:pong)(?=.*oc_chef_authz:pong).*\].*Keygen:map\[(?=.*keys:10)(?=.*max:10)(?=.*max_workers:2)(?=.*cur_max_workers:2)(?=.*avail_workers:)(?=.*start_size:0).*\]\}/)
  end
end
