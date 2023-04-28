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

	// Our new client shouldn't be allowed to delete our new node.
	if err = client2.Nodes.Delete("acltest"); err == nil {
		fmt.Fprintln(os.Stderr, "Expected error when deleting node without acl permission")
	}

	// Fetch existing ACL for our test node
	acls, err := client.ACLs.Get("nodes", "acltest")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue fetching acls for node:", err)
	}

	// Modify the delete ACL to allow our test client to delete the node
	acl, ok := acls["delete"]
	if !ok {
		fmt.Fprintln(os.Stderr, "Expected delete acl for node")
	}
	acl.Actors = []string{}
	acl.Clients = append(acl.Clients, "acltest")
	update := chef.ACL{"delete": acl}

	if err = client.ACLs.Put("nodes", "acltest", "delete", &update); err != nil {
		fmt.Fprintln(os.Stderr, "Issue updating acl for node:", err)
	}

	// Our new client should now be allowed to delete our new node.
	if err = client2.Nodes.Delete("acltest"); err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting node after setting acl:", err)
	}

	// Clean up
	if err := client.Clients.Delete("acltest"); err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting client for acl:", err)
	}
}
