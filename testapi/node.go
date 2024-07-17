package testapi

import (
	"fmt"
	"os"

	"github.com/go-chef/chef"
)

// node exercise the chef api
func Node() {
	// Use the default test org
	client := Client(nil)
	fmt.Printf("Client settings %+v\n client auth %+v\n", client, client.Auth)
	version := "1.0"
	if len(os.Args) > 6 {
		version = os.Args[6]
	}

	// List initial nodes
	nodeList, err := client.Nodes.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list nodes: ", err)
	}
	fmt.Println("List initial nodes", nodeList)

	// Define a Node object
	node1 := chef.NewNode("node1" + version)
	node1.RunList = []string{"pwn"}
	node1.AutomaticAttributes = map[string]interface{}{
		"attr": "value",
	}
	fmt.Printf("Define node1 %+v\n", node1)

	// Delete node1 ignoring errors :)
	err = client.Nodes.Delete(node1.Name)

	// Create
	nodeResult, err := client.Nodes.Post(node1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't create node node1. ", err)
	}
	fmt.Println("Added node1", nodeResult)

	// List nodes
	nodeList, err = client.Nodes.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list nodes: ", err)
	}
	fmt.Println("List nodes after adding node1", nodeList)

	// Create a second time
	nodeResult, err = client.Nodes.Post(node1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't recreate node node1. ", err)
	}
	cerr, err := chef.ChefError(err)
	if cerr != nil {
		fmt.Fprintln(os.Stderr, "Couldn't recreate node node1. ", cerr.StatusCode())
	}
	fmt.Println("Added node1", nodeResult)

	// Read node1 information
	serverNode, err := client.Nodes.Get(node1.Name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get node: ", err)
	}
	fmt.Printf("Get node1 %+v\n", serverNode)

	// Check node1 exits
	err = client.Nodes.Head(node1.Name)
	fmt.Println("Head node node1:", err)

	// Check nothere exits
	err = client.Nodes.Head("nothere")
	fmt.Println("Head node nothere:", err)

	// update node
	node1.RunList = append(node1.RunList, "recipe[works]")
	node1.AutomaticAttributes = map[string]interface{}{}
	updateNode, err := client.Nodes.Put(node1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't update node: ", err)
	}
	fmt.Println("Update node1", updateNode)

	// Info after update
	serverNode, err = client.Nodes.Get(node1.Name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get node: ", err)
	}
	fmt.Printf("Get node1 after update %+v\n", serverNode)

	// Delete node ignoring errors :)
	err = client.Nodes.Delete(node1.Name)
	fmt.Println("Delete node1", err)

	// List nodes
	nodeList, err = client.Nodes.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list nodes: ", err)
	}
	fmt.Println("List nodes after cleanup", nodeList)
}
