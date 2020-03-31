//
// Test the go-chef/chef chef server api /_status endpoints against a live server
//
package main

import (
	"chefapi_test/testapi"
	"fmt"
)

// main Exercise the chef server api
func main() {
	// Create a client for access
	client := testapi.Client()

	status, err := client.Status.Get()
	if err != nil {
		fmt.Println("Issue getting status information", err)
	}
	fmt.Printf("List status: %+v", status)
}
