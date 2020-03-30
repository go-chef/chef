# Inspec tests for the status chef api go module
#

describe command('/go/src/chefapi_test/bin/status') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/List status: \{Status:pong Upstreams:map\[(?=.*chef_solr:pong)(?=.*chef_sql:pong)(?=.*chef_index:pong)(?=.*oc_chef_action:pong)(?=.*oc_chef_authz:pong).*\].*Keygen:map\[(?=.*keys:10)(?=.*max:10)(?=.*max_workers:1)(?=.*cur_max_workers:1)(?=.*inflight:0)(?=.*avail_workers:1)(?=.*start_size:0).*\]\}/) }
end
