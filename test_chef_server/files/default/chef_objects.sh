#!/bin/bash -x

# Add two users
chef-server-ctl user-create user1 user1 mid last user1.last@nordstrom.com dummuy1pass
chef-server-ctl user-create user2 user2 mid last user2.last@nordstrom.com dummuy1pass

# Add a user to an org
#./chef-server-ctl org-user-add ORG_NAME USER_NAME
# Add an admin user to an org
#./chef-server-ctl org-user-add ORG_NAME USER_NAME -a

# define the current node to the chef server
# chef-client

# Add a cookbook
knife cookbook upload
