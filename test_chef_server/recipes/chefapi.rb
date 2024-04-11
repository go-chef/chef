# recipe chef_tester::chefapi
#
#  sudo snap refresh --classic --channel=1.20/stable go

package 'git'

snap_package 'go' do
  options 'classic'
  channel 'stable'
  action :upgrade
end
