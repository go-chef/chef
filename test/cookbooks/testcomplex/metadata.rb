name 'testcomplex'
maintainer 'The Authors'
maintainer_email 'you@example.com'
license 'All Rights Reserved'
description 'Installs/Configures testcomplex'
version '1.2.3'
chef_version '>= 18.0'

supports 'ubuntu', '>= 20.04'
supports 'redhat'

depends 'lvm', '~> 6.1' # Needed for VG and LV management
depends 'vagrant', '>= 4.0.14'

gem     'json', '>1.0.0'
