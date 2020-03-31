directory '/etc/chef/accepted_licenses' do
  recursive true
end

# Workaround for https://github.com/chef-cookbooks/chef-server/issues/161
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

apt_update 'update'

execute 'apt upgrade' do
  command 'DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" dist-upgrade'
  ignore_failure true
end
