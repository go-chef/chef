unified_mode true if respond_to?(:unified_mode)

property :a_prop, String, default: 'default'

action :create do
  file ::File.join('/tmp', new_resource.a_prop)
end
