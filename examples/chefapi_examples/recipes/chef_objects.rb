recipes/chef_objects.rb# Add chef objects to the server for testing

execute 'Set the host name' do
  command 'hostname testhost'
end

append_if_no_line 'add hostname to /etc/hosts' do
  line '127.0.01 testhost'
  path '/etc/hosts'
end

# Create an organization

execute 'create test organization' do
  command '/opt/opscode/bin/chef-server-ctl org-create test test_org'
  not_if '/opt/opscode/bin/chef-server-ctl org-list |grep test'
end

execute 'get the ssl certificate for the chef server' do
  command 'knife ssl fetch'
  not_if { File.exist? '/root/.chef/trusted_certs/testhost' }
end

# Register this node with the server
directory '/etc/chef'
file '/etc/chef/client.rb' do
  content '
chef_server_url "https://testhost/organizations/test"
client_fork true
file_backup_path "/var/chef/backup"
file_cache_path "/var/chef/cache"
log_location "/var/log/chef/client.log"
nodename "testhost"
validation_client_name "pivotal"
validation_key "/etc/opscode/pivotal.pem"
trusted_certs_dir "/root/.chef/trusted_certs"
ohai.disabled_plugins = [:C,:Cloud,:Rackspace,:Eucalyptus,:Command,:DMI,:Erlang,:Groovy,:IpScopes,:Java,:Lua,:Mono,:NetworkListeners,:Passwd,:Perl,:PHP,:Python]
'
end

directory '/fixtures/bin' do
  recursive true
end

directory '/fixtures/chef/cookbooks' do
  recursive true
end

cookbook_file '/fixtures/bin/chef_objects.sh' do
  source 'chef_objects.sh'
  mode '0755'
end

remote_directory '/fixtures/chef/cb' do
  source 'cb'
end

directory '/var/log/chef' do
  recursive true
end

directory '/var/chef' do
recipes/chef_objects.rb  recursive true
end
