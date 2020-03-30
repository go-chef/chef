package main

import (
	"fmt"
	"os"
	"time"

	"chefapi_test/testapi"
	"github.com/go-chef/chef"
)

func main() {
	// Add nodes
	client := testapi.Client()
	addNodes(client)
	// Give the nodes time to end up in the search data bases.  An immediate search will show no nodes
	time.Sleep(10 * time.Second)

	// Stanard search
	res, err := client.Search.Exec("node", "name:node*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running Search.Exec() ", err)
	}
	fmt.Printf("List nodes from Exec query Total:%+v\n", res.Total)

	// Partial search
	part := make(map[string]interface{})
	part["name"] = []string{"name"}
	pres, err := client.Search.PartialExec("node", "*:*", part)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running Search.PartialExec()", err)
	}
	fmt.Printf("List nodes from partial search Total:%+v\n", pres.Total)

	// Clean up nodes
	deleteNodes(client)
}

func addNodes(client *chef.Client) {
	for i := 0; i < 1200; i++ {
		node := chef.NewNode("node" + fmt.Sprintf("%d", i))
		_, err := client.Nodes.Post(node)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Issue adding node", node, err)
		}
	}
	return
}

func deleteNodes(client *chef.Client) {
	for i := 0; i < 1200; i++ {
		err := client.Nodes.Delete("node" + fmt.Sprintf("%d", i))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Issue deleting node", err)
		}
	}
	return
}
