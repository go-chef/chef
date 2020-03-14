//
// Test the go-chef/chef chef server api /updated_since endpoints against a live server
//
package main

import (
	"fmt"
	"chefapi_test/testapi"
)


// main Exercise the chef server api
func main() {
        // Create a client for access
	client := testapi.Client()

	universe, err := client.UpdatedSince.Get(1)
	if err != nil {
		fmt.Println("Issue getting universe information", err)
	}
	fmt.Printf("List updated_since initial: %+v", universe)

	// Define a Node object
	node1 := chef.NewNode("node1")
	node1.RunList = []string{"pwn"}
	_, err := client.Nodes.Post(node1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't create node node1. ", err)
	}

	universe, err := client.UpdatedSince.Get(1)
	if err != nil {
		fmt.Println("Issue getting universe information", err)
	}
	fmt.Printf("List updated_since initial: %+v", universe)

	// Delete node1 ignoring errors :)
	_ = client.Nodes.Delete(node1.Name)
}
