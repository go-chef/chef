// Test the go-chef/chef chef server api /universe endpoints against a live server
package testapi

import (
	"fmt"
)

// universe exercise the chef server api
func Universe() {
	// Create a client for access
	client := Client()

	universe, err := client.Universe.Get()
	if err != nil {
		fmt.Println("Issue getting universe information", err)
	}
	fmt.Printf("List universe: %+v", universe)
}
