# Inspec tests for the vault chef api go module
#
describe command('/go/src/chefapi_test/bin/vault') do
  its('stderr') { should_not match(/error|no such file|cannot find|not used|undefined/) }
  its('stdout') { should match(/^List vaults before creation/) }
  its('stdout') { should match(/^Created testv secrets vault item &\{DataBagItem:[0-9a-fx]+ Keys:[0-9a-fx]+ Name:secrets Vault:testv VaultService:[0-9a-fx]+\}/) }
  its('stdout') { should match(/^Updated testv secrets vault item/) }
  its('stdout') { should match(/^List the vault items \[secrets\]/) }
  its('stdout') { should match(/^List vaults after creation \[\]/) }
  its('stdout') { should match(/^Get testv secrets vault item&\{DataBagItem:[0-9a-fx]+ Keys:[0-9a-fx]+ Name:secrets Vault:testv VaultService:[0-9a-fx]+\}/) }
  its('stdout') { should match(/^Show testv databag item map\[(?=.*id:secrets)(?=.*foo:map\[)(?=.*auth_tag:[\w=+])(?=.*version:3)(?=.*iv:\w+)(?=.*encrypted_data:[\w\=\+]+)(?=.*cipher:aes-256-gcm).*\]/) }
  its('stdout') { should match(/^Show testv keys &\{DataBagItem:[0-9a-fx]+ Name:secrets\}/) }
  its('stdout') { should match(/^Show testv keys &\{DataBagItem:[0-9a-fx]+ Name:secrets\}/) }
  its('stdout') { should match(/^Show testv keys databagitem map\[(?=.*admins:\[pivotal\])(?=.*clients:\[\])(?=.*id:secrets_keys)(?=.*mode:default)(?=.*pivotal:\w+)(?=.*search_query:\[\]).*\]/) }
  its('stdout') { should match(/^List decrypted initial vault item values &map\[(?=.*id:secrets)(?=.*foo:bar).*\]/) }
  its('stdout') { should match(/^Updated based on get of testv secrets vault item/) }
  its('stdout') { should match(/^Delete testv secrets vault item/) }
  its('stdout') { should match(/^List vaults after deletion/) }
  its('stdout') { should match(/^Data bag deleted &\{Name:testv JsonClass:Chef::DataBag ChefType:data_bag\}/) }
end
