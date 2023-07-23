package chef

// TODO: Add a call to retrive all ACLs.

import (
	"errors"
	"fmt"
)

const adminUser = "pivotal"

type ACLService struct {
	client *Client
}

// ACL represents the native Go version of the deserialized ACL type
// The key string will be one of:
//  create, delete, grant, read, update
//
// indicating the type of access granted to an accessor in the ACLitems lists
type ACL map[string]ACLitems

// ACLitems
//
// Newer versions of the Chef server split Actors into Users and Clients.
type ACLitems struct {
	Actors  ACLitem `json:"actors"`
	Clients ACLitem `json:"clients"`
	Groups  ACLitem `json:"groups"`
	Users   ACLitem `json:"users"`
}

// ACLitem
// A list of a specific type of accessor.  In ACLitems the group ACLitem consists of a list of groups.
type ACLitem []string

func NewACL(acltype string, actors ACLitem, groups ACLitem, users ACLitem, clients ACLitem) *ACL {
	return &ACL{
		acltype: *NewACLItems(actors, groups, users, clients),
	}
}

func NewACLItems(actors ACLitem, groups ACLitem, users ACLitem, clients ACLitem) *ACLitems {
	return &ACLitems{
		Actors:  actors,
		Clients: clients,
		Groups:  groups,
		Users:   users,
	}
}

// ACLAdminAccess
// Verify that pivotal is in the Users access list for each type of ACL access
func ACLAdminAccess(acl *ACL) (err error) {
	err = errors.New("pivotal was not in a Users access list")
	// For each type of access in the acl
	// Verify that "pivotal" is listed in the Users list
	for _, item := range *acl {
		// See if "pivotal" is in the list of Users
		foundUser := false
		for _, user := range item.Users {
			if user == adminUser {
				// done with this item in the acl
				foundUser = true
				break
			}
		}
		if foundUser {
			// this item checks out, verify the next
			continue
		} else {
			return err
		}
	}
	return nil
}

// Get gets an ACL from the Chef server.
//
// Warning: This API is not included in the Chef documentation and thus probably not officially supported.
// Package documentation is based on the `knife` source code and packet capture.
// It could be wrong or change in future Chef server updates.
//
// Subkind can be one of: clients, containers, cookbook_artifacts, cookbooks, data, environments, groups, nodes, roles, policies, policy_groups.
//
//*
// map[create:{Groups:[admins clients users] Actors:[] Users:[pivotal] Clients:[]}
//     delete:{Groups:[admins users] Actors:[] Users:[pivotal] Clients:[]}
//     grant:{Groups:[admins] Actors:[] Users:[pivotal] Clients:[]}
//     read:{Groups:[admins clients users] Actors:[] Users:[pivotal] Clients:[]}
//     update:{Groups:[admins users] Actors:[] Users:[pivotal] Clients:[]}]
//*
// Returns the ACL for multiple perms (create, read, update, delete, grant).
// Older versions of the Chef server only include ACLs for "groups" and "actors."
// If you're using a more recent version then the contents of "actors" is split up in "users" and "clients."
func (a *ACLService) Get(subkind string, name string) (acl ACL, err error) {
	url := fmt.Sprintf("%s/%s/_acl?detail=granular", subkind, name)
	err = a.client.magicRequestDecoder("GET", url, nil, &acl)
	return
}

// Put updates an ACL on the Chef server.
//
// Warning: This API is not included in the Chef documentation and thus probably not officially supported.
// Package documentation is based on the `knife` source code and packet capture.
// It could be wrong or change in future Chef server updates.
//
// To change an ACL you have to fetch it from the Chef server first, as it expects the PUT request to
// contain the same elements as the GET response. While the GET response returns ACLs for all perms, you have
// to update each one separately.
//
// On newer versions of the Chef server you may need to replace "actors" with an empty list. Looks like the
// actors list is only included for backwards compatibility but can't be in the PUT request.
func (a *ACLService) Put(subkind, name string, perm string, item *ACL) (err error) {
	url := fmt.Sprintf("%s/%s/_acl/%s", subkind, name, perm)
	err = ACLAdminAccess(item)
	if err != nil {
		return
	}
	body, err := JSONReader(item)
	if err != nil {
		return
	}

	err = a.client.magicRequestDecoder("PUT", url, body, nil)
	return
}
