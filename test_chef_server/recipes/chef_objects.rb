# Add chef objects to the server for testing

# Create an organization
execute 'create test organization' do
  command '/opt/opscode/bin/chef-server-ctl org-create test test_org'
  ignore_failure true
  not_if '/opt/opscode/bin/chef-server-ctl org-list | grep test'
end

execute 'get the ssl certificate for the chef server' do
  command 'knife ssl fetch'
  ignore_failure true
  not_if { ::File.exist? '/root/.chef/trusted_certs/localhost' }
end

# Register this node with the server
directory '/etc/chef'
file '/etc/chef/client.rb' do
  content '
chef_server_url "https://localhost/organizations/test"
client_fork true
file_backup_path "/var/chef/backup"
file_cache_path "/var/chef/cache"
log_location "/var/log/chef/client.log"
nodename "localhost"
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
  recursive true
end

file '/etc/opscode/required' do
  content "file '/tmp/required'"
  mode '0600'
end

replace_or_add 'rr turn it on' do
  path '/etc/opscode/chef-server.rb'
  pattern "required_recipe['enable']"
  line "required_recipe['enable'] =  true"
end

replace_or_add 'rr set path' do
  path '/etc/opscode/chef-server.rb'
  pattern "required_recipe['path']"
  line "required_recipe['path'] =  '/etc/opscode/required'"
end

execute 'Change the chef options' do
  command 'chef-server-ctl reconfigure'
end
