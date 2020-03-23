package main

import (
	"fmt"
	"os"

	"chefapi_test/testapi"
	"github.com/go-chef/chef"
)

func main() {
	// Use the default test org
	client := testapi.Client()

	// List initial nodes
	nodeList, err := client.Nodes.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list nodes: ", err)
	}
	fmt.Println("List initial nodes", nodeList)

	// Define a Node object
	node1 := chef.NewNode("node1")
	node1.RunList = []string{"pwn"}
	fmt.Println("Define node1", node1)

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
	fmt.Println("Added node1", nodeResult)

	// Read node1 information
	serverNode, err := client.Nodes.Get("node1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get node: ", err)
	}
	fmt.Printf("Get node1 %+v\n", serverNode)

	// update node
	node1.RunList = append(node1.RunList, "recipe[works]")
	updateNode, err := client.Nodes.Put(node1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't update node: ", err)
	}
	fmt.Println("Update node1", updateNode)

	// Info after update
	serverNode, err = client.Nodes.Get("node1")
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
