# recipe chef_tester::chefapi
#

package 'git'

package 'golang'

directory '/go/src/github.com/go-chef' do
  recursive true
end

#  TODO: allow for testing a branch or the version of the api this cookbook is embedded in
git '/go/src/github.com/go-chef/chef' do
  repository 'https://github.com/go-chef/chef.git'
end

remote_directory 'local_go' do
  files_backup false
  path        '/go'
  purge       false
  recursive   true
  source      'go'
end

fileutils '/go/src/chefapi_test/bin' do
  file_mode ['+x']
end
