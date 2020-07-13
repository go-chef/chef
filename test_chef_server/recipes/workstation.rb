#
# Install chefdk workstation
#
#

package 'libx11.dev'

package 'ntpdate'

execute 'Sync the time' do
  command 'ntpdate time.nist.gov'
end

file '/etc/chef/accepted_licenses/chef_workstation' do
  content "---
    id: chef-workstation
    name: Chef Workstation
    date_accepted: '2020-05-06T23:18:26+00:00'
    accepting_product: chef-workstation
    accepting_product_version: 2.0.0
    user: vagrant
    file_format: 1"
end

execute 'Install chef workstation' do
  command 'curl --silent --show-error https://omnitruck.chef.io/install.sh | sudo -E bash -s -- -c stable -P chef-workstation --chef-license accept'
  not_if 'test -x /opt/chef-workstation/bin/chef'
end
