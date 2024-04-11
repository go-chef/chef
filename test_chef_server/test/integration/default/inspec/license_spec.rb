# Inspec tests for the license chef api go module
#

describe command('/go/src/github.com/go-chef/chef/testapi/bin/license') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stderr') { should_not match(/Issue/) }
  its('stdout') { should match(%r{^List license: {LimitExceeded:false NodeLicense:25 NodeCount:0 UpgradeUrl:https://www.chef.io/pricing}}) }
end
