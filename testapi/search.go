package testapi

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-chef/chef"
)

// search exercise the chef api
func Search() {
	// Add nodes
	client := Client(nil)
	addNodes(client)
	// Give the nodes time to end up in all of the data bases.  An immediate search will show no nodes
	time.Sleep(10 * time.Second)

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
	query.Rows = 2
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

	// Run the query, JSON output
	jres, err := query.DoJSON(client)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running query ", err)
	}
	fmt.Printf("List nodes from query JSON format %+v\n", jres)

	// Get the next page of results
	fmt.Printf("Query after the call %+v\n", query)
	query.Start = query.Start + query.Rows
	res, err = query.Do(client)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running 2nd query ", err)
	}
	fmt.Printf("List 2nd set of nodes from query %+v\n", res)

	// Get the next page of results again, in JSON format
	query.Start = query.Start - query.Rows
	jres, err = query.DoJSON(client)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running 2nd query ", err)
	}
	fmt.Printf("List 2nd set of nodes from query JSON format %+v\n", jres)

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

	// You can also use the service to run a query JSON Format
	jres, err = client.Search.ExecJSON("node", "name:node1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running Search.Exec() ", err)
	}
	fmt.Printf("List nodes from Exec query JSON format %+v\n", jres)
	// dump out results back in json as an example
	fmt.Println("JSON output example")
	os.Stdout.WriteString("\n")

	jres, err = client.Search.ExecJSON("node", "name:*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running Search.Exec() ", err)
	}
	fmt.Printf("List nodes from all nodes Exec query JSON format %+v\n", jres)

	// Partial search
	part := make(map[string]interface{})
	part["name"] = []string{"name"}
	jpres, err := client.Search.PartialExecJSON("node", "*:*", part)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running Search.PartialExec()", err)
	}
	fmt.Printf("List nodes from partial search %+v\n", jpres)

	// Partial search JSON format
	part = make(map[string]interface{})
	part["name"] = []string{"name"}
	jpres, err = client.Search.PartialExecJSON("node", "*:*", part)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue running Search.PartialExec()", err)
	}
	fmt.Printf("List nodes from partial search JSON format %+v\n", jpres)
	for i, row := range jpres.Rows {
		fmt.Fprintf(os.Stdout, "Partial search JSON format row: %v rawjson: %v\n", i, string(row.Data))
	}

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
