package main

import (
	"encoding/json"
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
	// Give the nodes time to end up in all of the data bases.  An immediate search will show no nodes
	time.Sleep(10 * time.Second)

	// TODO: Search limit is hardcoded to 1000, figure out how to do paging and to set the limit

	// List Indexes
	indexes, err := client.Search.Indexes()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing indexes ", err)
	}
	fmt.Printf("List indexes %+v EndIndex\n", indexes)

	// build an invalid seach query
	query, err := client.Search.NewQuery("node", "name")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue building invalid query", err)
	}

	// build a seach query
	query, err = client.Search.NewQuery("node", "name:node*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue building query ", err)
	}
	fmt.Printf("List new query %+v\n", query)

	// Run the query
	res, err := query.Do(client)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running query ", err)
	}
	fmt.Printf("List nodes from query %+v\n", res)

	// You can also use the service to run a query
	res, err = client.Search.Exec("node", "name:node1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running Search.Exec() ", err)
	}
	fmt.Printf("List nodes from Exec query %+v\n", res)
	// dump out results back in json as an example
	fmt.Println("JSON output example")
	jsonData, err := json.MarshalIndent(res, "", "\t")
	os.Stdout.Write(jsonData)
	os.Stdout.WriteString("\n")

	res, err = client.Search.Exec("node", "name:*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running Search.Exec() ", err)
	}
	fmt.Printf("List nodes from all nodes Exec query %+v\n", res)

	// Partial search
	part := make(map[string]interface{})
	part["name"] = []string{"name"}
	pres, err := client.Search.PartialExec("node", "*:*", part)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running Search.PartialExec()", err)
	}
	fmt.Printf("List nodes from partial search %+v\n", pres)

	// Clean up nodes
	deleteNodes(client)
}

func addNodes(client *chef.Client) {
	for i := 0; i < 2; i++ {
		node := chef.NewNode("node" + fmt.Sprintf("%d", i))
		_, err := client.Nodes.Post(node)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Issue adding node", node, err)
		}
		bode := chef.NewNode("bode" + fmt.Sprintf("%d", i))
		_, err = client.Nodes.Post(bode)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Issue adding node", node, err)
		}
	}
	return
}

func deleteNodes(client *chef.Client) {
	for i := 0; i < 2; i++ {
		err := client.Nodes.Delete("node" + fmt.Sprintf("%d", i))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Issue deleting node", err)
		}
		err = client.Nodes.Delete("bode" + fmt.Sprintf("%d", i))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Issue deleting node", err)
		}
	}
	return
}
