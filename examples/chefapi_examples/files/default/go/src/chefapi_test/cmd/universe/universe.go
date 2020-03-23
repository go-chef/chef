//
// Test the go-chef/chef chef server api /universe endpoints against a live server
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

	universe, err := client.Universe.Get()
	if err != nil {
		fmt.Println("Issue getting universe information", err)
	}
	fmt.Printf("List universe: %+v", universe)
}
