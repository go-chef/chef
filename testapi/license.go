//
// Test the go-chef/chef chef server api /license endpoints against a live server
//
package testapi

import (
	"fmt"
)

// license exercise the chef server api
func License() {
	// Create a client for access
	client := Client()

	license, err := client.License.Get()
	if err != nil {
		fmt.Println("Issue getting license information", err)
	}
	fmt.Printf("List license: %+v", license)
}
