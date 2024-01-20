# recipe chef_tester::chefapi
#

package 'git'

package 'golang-1.21'

link '/usr/bin/go' do
  link_type :symbolic
  to '/usr/lib/go-1.21/bin/go'
end
