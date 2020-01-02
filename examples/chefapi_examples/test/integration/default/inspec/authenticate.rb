# Inspec tests for the authenticate_user chef api go module
#

describe command('/go/src/chefapi_test/bin/authenticate') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/^Authenticate with a valid password \<nil\>/) }
  its('stdout') { should match(/^Authenticate with an invalid password POST https:\/\/localhost\/authenticate_user: 401/) }
end
