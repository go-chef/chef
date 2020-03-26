//
// Test the go-chef/chef chef server api /license endpoints against a live server
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

	license, err := client.License.Get()
	if err != nil {
		fmt.Println("Issue getting license information", err)
	}
	fmt.Printf("List license: %+v", license)
}
