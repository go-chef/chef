# Inspec tests for the principal chef api go module
#

describe command('/go/src/testapi/bin/principal') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  # client and user with the same name
  its('stdout') { should match(/Client principal \{Principals:\[{Name:client1 Type:user.*PublicKey:.*AuthzId:.*OrgMember:false.*\{Name:client1 Type:client.*OrgMember:true.*\}/) }
  its('stdout') { should match(/User principal \{Principals:\[\{Name:usr1 Type:user.*PublicKey:.*AuthzId:.*OrgMember:false.*\}/) }
end
