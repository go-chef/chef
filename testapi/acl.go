//
// Test the go-chef/chef chef server api ACL endpoints against a live server
//
package testapi

import (
	"fmt"
	"os"

	"github.com/go-chef/chef"
)

// ACL exercise the chef server api
func ACL() {
	client := Client()

	// Create a node

	node := chef.NewNode("acltest")
	_, err := client.Nodes.Post(node)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue adding node for acl:", err)
	}

	// Create a new client and another chef.Client that uses its private key

	newClient := chef.ApiNewClient{
		Name:       "acltest",
		ClientName: "acltest",
		CreateKey:  true,
	}
	aclClient, err := client.Clients.Create(newClient)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue adding client for acl:", err)
	}

	// We want exactly the same test API client but with a different key.

	client2 := Client()
	client2.Auth.ClientName = "acltest"
	private, err := chef.PrivateKeyFromString([]byte(aclClient.ChefKey.PrivateKey))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue creating private key from client create response:", err)
	}
	client2.Auth.PrivateKey = private

	// Test: Our new client isn't allowed to delete our new node.

	if err = client2.Nodes.Delete("acltest"); err == nil {

		fmt.Fprintln(os.Stderr, "Issue expected error when deleting node without acl permission")
	}

	// Test Fetch existing ACL for our test node

	acls, err := client.ACLs.Get("nodes", "acltest")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue fetching acls for node:", err)
	}
	// TODO: Test the values returned

	// Test: Modify the delete ACL to allow our test client to delete the node

	acl, ok := acls["delete"]
	if !ok {
		fmt.Fprintln(os.Stderr, "Issue expected a delete acl list for the node")
	}
	acl.Actors = []string{}
	acl.Clients = append(acl.Clients, "acltest")
	update := chef.ACL{"delete": acl}
	if err = client.ACLs.Put("nodes", "acltest", "delete", &update); err != nil {
		fmt.Fprintln(os.Stderr, "Issue updating acl for node:", err)
	}

	// Test: Our new client should now be allowed to delete our new node.

	if err = client2.Nodes.Delete("acltest"); err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting node after setting acl:", err)
	}

	// Test: Verify that the admin credentials are present in the fetched acl

	err = chef.ACLAdminAccess(&update)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue checking for pivtoal user in acl:", err)
	}

	// Test: Verify that the admin credentials are present in the update acl

	err = chef.ACLAdminAccess(&update)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue checking for pivtoal user in acl:", err)
	}

	// Test: Remove the pivotal and and verify that the admin credentials are not present in the acl

	acl, ok = acls["delete"]
	acl.Users = []string{}
	update2 := chef.ACL{"create": acls["create"], "delete": acl, "grant": acls["grant"], "read": acls["read"], "update": acls["update"]}
	err = chef.ACLAdminAccess(&update2)
	if err == nil {
		fmt.Fprintln(os.Stderr, "Issue expected missing user checking for pivotal user in acl:", err)
	}

	// Test: Try to update the acl without the pivotal user

	if err = client.ACLs.Put("nodes", "acltest", "delete", &update2); err == nil {
		fmt.Fprintln(os.Stderr, "Issue expected missing user credentials  updating acl for node missing pivotal:", err)
	}

	// TODO: Fetch fail tests
	// TODO: Wrong class failures

	// Clean up
	if err := client.Clients.Delete("acltest"); err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting client for acl:", err)
	}
}
