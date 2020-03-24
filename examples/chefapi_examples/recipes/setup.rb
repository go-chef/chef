directory '/etc/chef/accepted_licenses' do
  recursive true
end

file '/etc/chef/accepted_licenses/chef_infra_client' do
  content "---
id: infra-client
name: Chef Infra Client
date_accepted: '2019-10-20T14:46:27+00:00'
accepting_product: infra-server
accepting_product_version: 0.6.0
user: vagrant
file_format: 1"
end

file '/etc/chef/accepted_licenses/chef_infra_server' do
  content "---
 id: infra-server
 name: Chef Infra Server
 date_accepted: '2019-10-20T14:46:27+00:00'
 accepting_product: infra-server
 accepting_product_version: 0.6.0
 user: vagrant
 file_format: 1"
end

file '/etc/chef/accepted_licenses/inspec' do
  content "---
id: inspec
name: Chef InSpec
date_accepted: '2019-10-20T14:46:27+00:00'
accepting_product: infra-server
accepting_product_version: 0.6.0
user: vagrant
file_format: 1"
end

apt_update 'update'

execute 'apt upgrade' do
  command 'DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" dist-upgrade'
  ignore_failure true
end
