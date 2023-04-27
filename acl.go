package chef

import "fmt"

type ACLService struct {
	client *Client
}

// ACL represents the native Go version of the deserialized ACL type
type ACL map[string]ACLitems

// ACLitems
//
// Newer versions of the Chef server split Actors into Users and Clients.
type ACLitems struct {
	Groups  ACLitem `json:"groups"`
	Actors  ACLitem `json:"actors"`
	Users   ACLitem `json:"users"`
	Clients ACLitem `json:"clients"`
}

// ACLitem
type ACLitem []string

func NewACL(acltype string, actors, groups ACLitem, users ACLitem, clients ACLitem) (acl *ACL) {
	acl = &ACL{
		acltype: ACLitems{
			Actors: actors,
			Groups: groups,
			Users: users,
			Clients: clients,
		},
	}
	return
}

// Get gets an ACL from the Chef server.
//
// Warning: This API is not included in the Chef documentation and thus probably not officially supported.
// Package documentation is based on the `knife` source code and packet capture.
// It could be wrong or change in future Chef server updates.
//
// Subkind can be one of: clients, containers, cookbook_artifacts, cookbooks, data, environments, groups, nodes, roles, policies, policy_groups.
//
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
	body, err := JSONReader(item)
	if err != nil {
		return
	}

	err = a.client.magicRequestDecoder("PUT", url, body, nil)
	return
}
