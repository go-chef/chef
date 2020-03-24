#
# Cookbook:: chefapi_examples
# Recipe:: default
#
# Copyright:: 2019, The Authors, All Rights Reserved.

hostname 'testhost'

node.override['chef-server']['api_fqdn'] = 'testhost'
include_recipe 'chef-server'
