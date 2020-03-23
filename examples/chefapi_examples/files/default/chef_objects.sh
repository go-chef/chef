#!/bin/bash -x

# Add two users
chef-server-ctl user-create user1 user1 mid last user1.last@nordstrom.com dummuy1pass
chef-server-ctl user-create user2 user2 mid last user2.last@nordstrom.com dummuy1pass

# Add a user to an org
./chef-server-ctl org-user-add ORG_NAME USER_NAME
# Add an admin user to an org
./chef-server-ctl org-user-add ORG_NAME USER_NAME -a

# define the current node to the chef server
chef-client

# Add a cookbook
knife cookbook upload

# Add a role

# Add the user to a group

# Add an environment

# Add a data bag

# Add a chef vault item

# Add a node and set the run list

# tag a node
