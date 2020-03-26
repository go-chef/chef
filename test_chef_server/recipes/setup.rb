apt_update 'update'

execute 'apt upgrade' do
  command 'DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" dist-upgrade'
  ignore_failure true
end
