//
// Test the go-chef/chef chef server api /license endpoints against a live server
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
	
	license, err := client.License.Get()
	if err != nil {
		fmt.Println("Issue getting license information", err)
	}
	fmt.Printf("List license: %+v", license)
}
