# Inspec tests for the universe chef api go module
#

describe command('/go/src/github.com/go-chef/chef/testapi/bin/universe') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(%r{^List universe: \{Books:map\[sampbook:\{Versions:map\[(?=.*0.2.0:\{LocationPath:https:\/\/testhost\/organizations\/test\/cookbooks\/sampbook\/0.2.0 LocationType:chef_server Dependencies:map\[\]\})(?=.*0.1.0:\{LocationPath:https:\/\/testhost\/organizations\/test\/cookbooks\/sampbook\/0.1.0 LocationType:chef_server Dependencies:map\[\]\}).*\]\}\]\}}) }
end
