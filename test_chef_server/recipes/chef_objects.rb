# Add chef objects to the server for testing

# Create an organization
execute 'create test organization' do
  command '/opt/opscode/bin/chef-server-ctl org-create test test_org'
  not_if '/opt/opscode/bin/chef-server-ctl org-list |grep test'
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
